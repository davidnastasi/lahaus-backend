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

type SignInSuite struct {
	suite.Suite
	mockCtrl      *gomock.Controller
	database      *mocks.MockStorageManager
	signInUseCase *users.SignInUserUseCase
}

func TestSignInSuite(t *testing.T) {
	suite.Run(t, new(SignInSuite))
}

func (suite *SignInSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.database = mocks.NewMockStorageManager(suite.mockCtrl)
	suite.signInUseCase = users.NewSignInUserUseCase(suite.database)
}

func (suite *SignInSuite) TearDownSuite() {
	suite.mockCtrl.Finish()
}

func (suite *SignInSuite) TestSignInUseCase_ExecuteSuccess() {
	user := &model.User{Password: "1"}
	suite.database.EXPECT().SaveUser(gomock.Any()).Return(nil)
	err := suite.signInUseCase.Execute(user)
	suite.NoError(err)
	suite.Equal("a4ayc_80_OGda4BO_1o_V0etpOqiLx1JwB5S3beHW0s=", user.Password)
}

func (suite *SignInSuite) TestSignInUseCase_ExecuteError() {
	suite.database.EXPECT().SaveUser(gomock.Any()).Return(errors.New("failed to get user"))
	err := suite.signInUseCase.Execute(&model.User{})
	suite.Error(err)
}
