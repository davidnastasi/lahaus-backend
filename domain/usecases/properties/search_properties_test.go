package properties_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"lahaus/domain/model"

	"lahaus/domain/usecases/properties"
	"lahaus/domain/usecases/properties/mocks"
	"testing"
)

type SearchPropertySuite struct {
	suite.Suite
	mockCtrl      *gomock.Controller
	database      *mocks.MockStorageManager
	searchUseCase *properties.SearchPropertyUseCase
}

func TestSearchPropertySuite(t *testing.T) {
	suite.Run(t, new(SearchPropertySuite))
}

func (suite *SearchPropertySuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.database = mocks.NewMockStorageManager(suite.mockCtrl)
	suite.searchUseCase = properties.NewSearchPropertyUseCase(suite.database)
}

func (suite *SearchPropertySuite) TearDownSuite() {
	suite.mockCtrl.Finish()
}

func (suite *SearchPropertySuite) TestSearchPropertyUseCase_ExecuteSuccess() {
	suite.database.EXPECT().FilterProperties(gomock.Any()).Return(&model.PropertiesPaging{
		Page:       1,
		PageSize:   10,
		TotalPages: 0,
		Total:      1,
		Data:       nil,
	}, nil)
	propertiesResult, err := suite.searchUseCase.Execute(properties.PropertySearchParams{})
	suite.NoError(err)
	suite.Equal(int64(1), propertiesResult.TotalPages)
}

func (suite *SearchPropertySuite) TestSearchPropertyUseCase_ExecuteError() {
	suite.database.EXPECT().FilterProperties(gomock.Any()).Return(nil, errors.New("fail to save in database"))
	propertiesResult, err := suite.searchUseCase.Execute(properties.PropertySearchParams{})
	suite.Error(err)
	suite.Nil(propertiesResult)
}
