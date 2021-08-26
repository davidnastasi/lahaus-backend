package users_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"lahaus/domain/model"
	"lahaus/domain/usecases/users"
	"lahaus/domain/usecases/users/mocks"
	"testing"
)

type ListFavouritesSuite struct {
	suite.Suite
	mockCtrl              *gomock.Controller
	database              *mocks.MockStorageManager
	listFavouritesUseCase *users.ListFavouritesUseCase
}

func TestListFavouriteSuite(t *testing.T) {
	suite.Run(t, new(ListFavouritesSuite))
}

func (suite *ListFavouritesSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.database = mocks.NewMockStorageManager(suite.mockCtrl)
	suite.listFavouritesUseCase = users.NewListFavouriteUseCase(suite.database)
}

func (suite *ListFavouritesSuite) TearDownSuite() {
	suite.mockCtrl.Finish()
}

func (suite *ListFavouritesSuite) TestListFavouritesUseCase_ExecuteSuccess() {
	suite.database.EXPECT().ListFavourites(gomock.Any()).Return(&model.PropertiesPaging{
		Page:       1,
		PageSize:   10,
		TotalPages: 0,
		Total:      1,
		Data:       nil,
	}, nil)
	propertiesResult, err := suite.listFavouritesUseCase.Execute(users.FavouritesSearchParams{})
	suite.NoError(err)
	suite.Equal(int64(1), propertiesResult.TotalPages)
}

func (suite *ListFavouritesSuite) TestListFavouritesUseCase_ExecuteError() {
	suite.database.EXPECT().ListFavourites(gomock.Any()).Return(nil, errors.New("fail to save in database"))
	propertiesResult, err := suite.listFavouritesUseCase.Execute(users.FavouritesSearchParams{})
	suite.Error(err)
	suite.Nil(propertiesResult)
}
