package api

import (
	"context"
	"errors"
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
)

type UserSuite struct {
	suite.Suite
	mockCtrl       *gomock.Controller
	userHandler    *UserHandler
	signUpExecutor *mocks.MockSignUpUserExecutor
	signInExecutor *mocks.MockSignInUserExecutor
	listExecutor   *mocks.MockListFavouritesExecutor
	addExecutor    *mocks.MockAddFavouriteExecutor
	chiRouter      *chi.Mux
	httpTest       *httptest.Server
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (suite *UserSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.signUpExecutor = mocks.NewMockSignUpUserExecutor(suite.mockCtrl)
	suite.signInExecutor = mocks.NewMockSignInUserExecutor(suite.mockCtrl)
	suite.listExecutor = mocks.NewMockListFavouritesExecutor(suite.mockCtrl)
	suite.addExecutor = mocks.NewMockAddFavouriteExecutor(suite.mockCtrl)
	suite.userHandler = NewUserHandler(suite.signInExecutor, suite.signUpExecutor, suite.addExecutor, suite.listExecutor)

	suite.chiRouter = chi.NewRouter()
	suite.chiRouter.Use(middleware.RequestID)
	suite.chiRouter.Route("/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", suite.userHandler.SignInUser)
			r.Post("/login", suite.userHandler.SignUpUser)
			r.Route("/me/favourites", func(r chi.Router) {
				r.Post("/", suite.userHandler.AddFavourite)
				r.Get("/", suite.userHandler.ListFavourites)
			})
		})
	})
	suite.httpTest = httptest.NewServer(suite.chiRouter)
}

func (suite *UserSuite) TearDownSuite() {
	suite.httpTest.Close()
	suite.mockCtrl.Finish()
}

func (suite *UserSuite) TestSignInUser_BadRequest() {
	req, err := http.NewRequest("POST", "/v1/users", strings.NewReader(`
		{
			"email": "code-challenge-lahaus@test.lh
		}
	`))
	suite.NoError(err)
	rr := httptest.NewRecorder()
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *UserSuite) TestSignInUser_InvalidEmail() {
	req, err := http.NewRequest("POST", "/v1/users", strings.NewReader(`
		{
			"email": "code-challenge-lahaustest.lh"
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *UserSuite) TestSignInUser_Success() {
	req, err := http.NewRequest("POST", "/", strings.NewReader(`
		{
			"email": "code-challenge-lahaus@test.lh"
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()

	suite.signInExecutor.EXPECT().Execute(gomock.Any()).Return(nil)
	handler := http.HandlerFunc(suite.userHandler.SignInUser)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test")
	handler.ServeHTTP(rr, req.WithContext(ctx))

	suite.Equal(http.StatusCreated, rr.Code)
}

func (suite *UserSuite) TestSignInUser_Error() {
	req, err := http.NewRequest("POST", "/v1/users", strings.NewReader(`
		{
			"email": "code-challenge-lahaus@test.lh"
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()
	suite.signInExecutor.EXPECT().Execute(gomock.Any()).Return(errors.New("fail to save"))
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusInternalServerError, rr.Code)
}

func (suite *UserSuite) TestSignUpUser_Error() {
	req, err := http.NewRequest("POST", "/v1/users/login", strings.NewReader(`
		{
			"email": "code-challenge-lahaus@test.lh",
			"password": "code-challenge-lahaus@test.lh"
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()

	suite.signUpExecutor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return("", errors.New("fail to get"))
	suite.chiRouter.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)
}

func (suite *UserSuite) TestSignUpUser_Success() {
	req, err := http.NewRequest("POST", "/v1/users/login", strings.NewReader(`
		{
			"email": "code-challenge-lahaus@test.lh",
			"password": "code-challenge-lahaus@test.lh"
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()

	suite.signUpExecutor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return("token", nil)
	suite.chiRouter.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
}

func (suite *UserSuite) TestSignUpUser_UserNotFound() {
	req, err := http.NewRequest("POST", "/v1/users/login", strings.NewReader(`
		{
			"email": "code-challenge-lahaus@test",
			"password": "code-challenge-lahaus@test"
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()

	suite.signUpExecutor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return("", model.NewUnauthorizedError(errors.New("user not found")))
	suite.chiRouter.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)
}

func (suite *UserSuite) TestSignUpUser_IncorrectPassword() {
	req, err := http.NewRequest("POST", "/v1/users/login", strings.NewReader(`
		{
			"email": "code-challenge-lahaus@test",
			"password": "code-challenge-lahaus@test"
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()

	suite.signUpExecutor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return("", model.NewUnauthorizedError(errors.New("invalid credentials")))
	suite.chiRouter.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)
}

func (suite *UserSuite) TestAddFavourites_BadRequest() {
	req, err := http.NewRequest("POST", "/v1/users/me/favourites", strings.NewReader(`
		{
			"propertyId":"1"
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()

	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *UserSuite) TestAddFavourite_Success() {
	req, err := http.NewRequest("POST", "/v1/users/me/favourites", strings.NewReader(`
		{
			"propertyId":1
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()

	suite.addExecutor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil)
	handler := http.HandlerFunc(suite.userHandler.AddFavourite)
	ctx := context.WithValue(req.Context(), "user", map[string]interface{}{
		"email":  "nn@nn.com",
		"userId": int64(1),
	})
	req = req.WithContext(ctx)
	ctxMid := context.WithValue(req.Context(), middleware.RequestIDKey, "test")
	req = req.WithContext(ctxMid)
	handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusCreated, rr.Code)
}

func (suite *UserSuite) TestAddFavourite_PropertyNotFound() {
	req, err := http.NewRequest("POST", "/v1/users/me/favourites", strings.NewReader(`
		{
			"propertyId":1
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()

	suite.addExecutor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(model.NewEntityNotFoundError(errors.New("fail to save favourite")))
	ctx := context.WithValue(req.Context(), "user", map[string]interface{}{
		"email":  "nn@nn.com",
		"userId": int64(1),
	})
	req = req.WithContext(ctx)
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusNotFound, rr.Code)
}

func (suite *UserSuite) TestAddFavourite_FailToSave() {
	req, err := http.NewRequest("POST", "/v1/users/me/favourites", strings.NewReader(`
		{
			"propertyId":1
		}
	`))
	suite.NoError(err)

	rr := httptest.NewRecorder()

	suite.addExecutor.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(errors.New("fail to save favourite"))
	ctx := context.WithValue(req.Context(), "user", map[string]interface{}{
		"email":  "nn@nn.com",
		"userId": int64(1),
	})
	req = req.WithContext(ctx)
	ctxMid := context.WithValue(req.Context(), middleware.RequestIDKey, "test")
	req = req.WithContext(ctxMid)
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusInternalServerError, rr.Code)
}

func (suite *UserSuite) TestListFavourites_IncorrectQueryParams() {
	req, err := http.NewRequest("GET", "/v1/users/me/favourites/?page=A", nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "user", map[string]interface{}{
		"email":  "nn@nn.com",
		"userId": int64(1),
	})
	req = req.WithContext(ctx)
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusBadRequest, rr.Code)
}

func (suite *UserSuite) TestListFavourites_FailToGet() {
	req, err := http.NewRequest("GET", "/v1/users/me/favourites/", nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()

	suite.listExecutor.EXPECT().Execute(gomock.Any()).Return(nil, errors.New("fail to get favourites"))
	ctx := context.WithValue(req.Context(), "user", map[string]interface{}{
		"email":  "nn@nn.com",
		"userId": int64(1),
	})
	req = req.WithContext(ctx)
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusInternalServerError, rr.Code)
}

func (suite *UserSuite) TestListFavourites_Success() {
	req, err := http.NewRequest("GET", "/v1/users/me/favourites/", nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()

	suite.listExecutor.EXPECT().Execute(gomock.Any()).Return(&model.PropertiesPaging{}, nil)
	ctx := context.WithValue(req.Context(), "user", map[string]interface{}{
		"email":  "nn@nn.com",
		"userId": int64(1),
	})
	req = req.WithContext(ctx)
	suite.chiRouter.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
}
