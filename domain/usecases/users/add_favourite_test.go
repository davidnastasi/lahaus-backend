package users_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"lahaus/domain/usecases/users"
	"lahaus/domain/usecases/users/mocks"
	"testing"
)

type AddFavouriteSuite struct {
	suite.Suite
	mockCtrl            *gomock.Controller
	database            *mocks.MockStorageManager
	addFavouriteUseCase *users.AddFavouriteUseCase
}

func TestAddFavouriteSuite(t *testing.T) {
	suite.Run(t, new(AddFavouriteSuite))
}

func (suite *AddFavouriteSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.database = mocks.NewMockStorageManager(suite.mockCtrl)
	suite.addFavouriteUseCase = users.NewAddFavouriteUseCase(suite.database)
}

func (suite *AddFavouriteSuite) TearDownSuite() {
	suite.mockCtrl.Finish()
}

func (suite *AddFavouriteSuite) TestAddFavouriteUseCase_ExecuteSuccess() {
	suite.database.EXPECT().GetProperty(gomock.Any()).Return(nil, true, nil)
	suite.database.EXPECT().AddFavourite(gomock.Any(), gomock.Any()).Return(nil)
	err := suite.addFavouriteUseCase.Execute(1, 1)
	suite.NoError(err)
}

func (suite *AddFavouriteSuite) TestAddFavouriteUseCase_ExecuteError_GetProperty() {

	suite.database.EXPECT().GetProperty(gomock.Any()).Return(nil, false, errors.New("fail"))
	err := suite.addFavouriteUseCase.Execute(1, 1)
	suite.Error(err)
}

func (suite *AddFavouriteSuite) TestAddFavouriteUseCase_ExecuteError_GetPropertyNotFound() {
	suite.database.EXPECT().GetProperty(gomock.Any()).Return(nil, false, nil)
	err := suite.addFavouriteUseCase.Execute(1, 1)
	suite.Error(err)
}

func (suite *AddFavouriteSuite) TestAddFavouriteUseCase_ExecuteError_AddFavouriteError() {
	suite.database.EXPECT().GetProperty(gomock.Any()).Return(nil, true, nil)
	suite.database.EXPECT().AddFavourite(gomock.Any(), gomock.Any()).Return(errors.New("fail"))
	err := suite.addFavouriteUseCase.Execute(1, 1)
	suite.Error(err)
}
