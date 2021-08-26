package api

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"lahaus/domain/model"
	"lahaus/infrastructure/api/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type PropertySuite struct {
	suite.Suite
	mockCtrl               *gomock.Controller
	propertyCreateExecutor *mocks.MockPropertyExecutor
	propertyUpdateExecutor *mocks.MockPropertyExecutor
	propertySearchExecutor *mocks.MockSearchPropertyExecutor
	propertyHandler        *PropertyHandler
	chiRouter              *chi.Mux
	httpTest               *httptest.Server
}

func TestPropertySuite(t *testing.T) {
	suite.Run(t, new(PropertySuite))
}

func (suite *PropertySuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.propertyCreateExecutor = mocks.NewMockPropertyExecutor(suite.mockCtrl)
	suite.propertyUpdateExecutor = mocks.NewMockPropertyExecutor(suite.mockCtrl)
	suite.propertySearchExecutor = mocks.NewMockSearchPropertyExecutor(suite.mockCtrl)
	suite.propertyHandler = NewPropertyHandler(suite.propertyCreateExecutor, suite.propertyUpdateExecutor, suite.propertySearchExecutor)

	suite.chiRouter = chi.NewRouter()
	suite.chiRouter.Use(middleware.RequestID)
	suite.chiRouter.Route("/v1", func(r chi.Router) {
		r.Route("/properties", func(r chi.Router) {
			r.Post("/", suite.propertyHandler.CreateProperty)
			r.Put("/{id}", suite.propertyHandler.UpdateProperty)
			r.Get("/", suite.propertyHandler.SearchProperties)
		})
	})

	suite.httpTest = httptest.NewServer(suite.chiRouter)
}

func (suite *PropertySuite) TearDownSuite() {
	suite.httpTest.Close()
	suite.mockCtrl.Finish()
}

func (suite *PropertySuite) TestCreateProperty_PropertyTypeError() {
	req, err := http.NewRequest("POST", "/v1/properties/", strings.NewReader(`
		{
			"title": "Apartamento cerca a la estación",
			"location": {
				"longitude": -94.0665887,
				"latitude": 94.6371593
			},
			"pricing": {
				"salePrice": 450000000
			},
			"propertyType": "UNKNOWN",
			"bedrooms": 3,
			"bathrooms": 2,
			"parkingSpots": 1,
			"area": 60,
			"photos": [
				"https://cdn.pixabay.com/photo/2014/08/11/21/39/wall-416060_960_720.jpg",
				"https://cdn.pixabay.com/photo/2016/09/22/11/55/kitchen-1687121_960_720.jpg"
			]
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *PropertySuite) TestCreateProperty_Error() {
	req, err := http.NewRequest("POST", "/v1/properties/", strings.NewReader(`
		{
			"title": "Apartamento cerca a la estación",
			"location": {
				"longitude": -94.0665887,
				"latitude": 94.6371593
			},
			"pricing": {
				"salePrice": 450000000
			},
			"propertyType": "HOUSE",
			"bedrooms": 3,
			"bathrooms": 2,
			"parkingSpots": 1,
			"area": 60,
			"photos": [
				"https://cdn.pixabay.com/photo/2014/08/11/21/39/wall-416060_960_720.jpg",
				"https://cdn.pixabay.com/photo/2016/09/22/11/55/kitchen-1687121_960_720.jpg"
			]
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.propertyCreateExecutor.EXPECT().Execute(gomock.Any()).Return(nil, errors.New("fail to save"))
	suite.chiRouter.ServeHTTP(rr, req)

	suite.Equal(http.StatusInternalServerError, rr.Code)
}

func (suite *PropertySuite) TestCreateProperty_Success() {
	req, err := http.NewRequest("POST", "/v1/properties/", strings.NewReader(`
		{
			"title": "Apartamento cerca a la estación",
			"location": {
				"longitude": -94.0665887,
				"latitude": 4.6371593
			},
			"pricing": {
				"salePrice": 450000000
			},
			"propertyType": "HOUSE",
			"bedrooms": 3,
			"bathrooms": 2,
			"area": 60,
			"photos": [
				"https://cdn.pixabay.com/photo/2014/08/11/21/39/wall-416060_960_720.jpg",
				"https://cdn.pixabay.com/photo/2016/09/22/11/55/kitchen-1687121_960_720.jpg"
			]
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()
	propertySaved := &model.Property{
		ID:    1,
		Title: "Apartamento cerca a la estación",
		Location: model.Location{
			Longitude: -94.0665887,
			Latitude:  4.6371593,
		},
		Pricing: model.Pricing{
			SalePrice: 450000000,
		},
		PropertyType: model.HOUSE,
		Bedrooms:     3,
		Bathrooms:    2,
		ParkingSpots: nil,
		Area:         60,
		Photos: model.Photos{
			"https://cdn.pixabay.com/photo/2014/08/11/21/39/wall-416060_960_720.jpg",
			"https://cdn.pixabay.com/photo/2016/09/22/11/55/kitchen-1687121_960_720.jpg",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    model.INACTIVE,
	}
	suite.propertyCreateExecutor.EXPECT().Execute(gomock.Any()).Return(propertySaved, nil)
	suite.chiRouter.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	js, err := simplejson.NewJson(rr.Body.Bytes())
	suite.NoError(err)

	id, err := js.Get("id").Int64()
	suite.NoError(err)
	suite.Equal(propertySaved.ID, id)
	_, found := js.CheckGet("administrativeFee")
	suite.False(found)
	_, found = js.CheckGet("description")
	suite.False(found)
	_, found = js.CheckGet("parkingSpots")
	suite.False(found)
	status, err := js.Get("status").String()
	suite.NoError(err)
	suite.Equal(string(propertySaved.Status), status)
}

func (suite *PropertySuite) TestUpdateProperty_InvalidParam() {
	req, err := http.NewRequest("PUT", "/v1/properties/A", strings.NewReader(`
		{
			"title": "Apartamento cerca a la estación",
			"description": "Apartamento cerca a la estación",
			"location": {
				"longitude": -94.0665887,
				"latitude": 4.6371593
			},
			"pricing": {
				"salePrice": 450000000
			},
			"propertyType": "HOUSE",
			"bedrooms": 3,
			"bathrooms": 2,
			"area": 60,
			"photos": [
				"https://cdn.pixabay.com/photo/2014/08/11/21/39/wall-416060_960_720.jpg",
				"https://cdn.pixabay.com/photo/2016/09/22/11/55/kitchen-1687121_960_720.jpg"
			]
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *PropertySuite) TestUpdateProperty_BadRequest() {
	req, err := http.NewRequest("PUT", "/v1/properties/1", strings.NewReader(`
		{
			"title": "Apartamento cerca a la estación",
			"description": "Apartamento cerca a la estación",
			"location": {
				"longitude": -94.0665887,
				"latitude": 4.6371593
			},
			"pricing": {
				"salePrice": 450000000
			},
			"propertyType": "HOUSE",
			"bedrooms": 3,
			"bathrooms": 2,
			"area": 60
			"photos": [
				"https://cdn.pixabay.com/photo/2014/08/11/21/39/wall-416060_960_720.jpg",
				"https://cdn.pixabay.com/photo/2016/09/22/11/55/kitchen-1687121_960_720.jpg"
			]
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *PropertySuite) TestUpdateProperty_Success() {
	req, err := http.NewRequest("PUT", "/v1/properties/1", strings.NewReader(`
		{
			"title": "Apartamento cerca a la estación",
			"description": "Apartamento cerca a la estación",
			"location": {
				"longitude": -94.0665887,
				"latitude": 4.6371593
			},
			"pricing": {
				"salePrice": 450000000
			},
			"propertyType": "HOUSE",
			"bedrooms": 3,
			"bathrooms": 2,
			"area": 60,
			"photos": [
				"https://cdn.pixabay.com/photo/2014/08/11/21/39/wall-416060_960_720.jpg",
				"https://cdn.pixabay.com/photo/2016/09/22/11/55/kitchen-1687121_960_720.jpg"
			]
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.propertyUpdateExecutor.EXPECT().Execute(gomock.Any()).Return(&model.Property{}, nil)
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
}

func (suite *PropertySuite) TestListProperty_BadRequestStatus() {
	req, err := http.NewRequest("GET", "/v1/properties/?status=A", nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *PropertySuite) TestListProperty_BadRequestLocationSize() {
	req, err := http.NewRequest("GET", "/v1/properties/?status=ACTIVE&bbox=1", nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *PropertySuite) TestListProperty_BadRequestInvalidPage() {
	req, err := http.NewRequest("GET", "/v1/properties/?status=ACTIVE&bbox=-1,1,-1,1&page=A", nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *PropertySuite) TestListProperty_BadRequestInvalidPageSize() {
	req, err := http.NewRequest("GET", "/v1/properties/?status=ACTIVE&bbox=-1,1,-1,1&page=1&pageSize=A", nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *PropertySuite) TestListProperty_BadRequestOutOfBoundPageSize() {
	req, err := http.NewRequest("GET", "/v1/properties/?status=ACTIVE&bbox=-1,1,-1,1&page=1&pageSize=30", nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *PropertySuite) TestListProperty_Error() {
	req, err := http.NewRequest("GET", "/v1/properties/?status=ACTIVE&bbox=-1,1,-1,1&page=1&pageSize=15", nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.propertySearchExecutor.EXPECT().Execute(gomock.Any()).Return(nil, errors.New("error fetching database"))
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *PropertySuite) TestListProperty_Success() {
	req, err := http.NewRequest("GET", "/v1/properties/?status=ACTIVE&bbox=-1,1,-1,1&page=1&pageSize=15", nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.propertySearchExecutor.EXPECT().Execute(gomock.Any()).Return(&model.PropertiesPaging{}, nil)
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
}
