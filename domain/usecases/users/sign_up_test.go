package users_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"lahaus/config"
	"lahaus/domain/model"
	"lahaus/domain/usecases/users"
	"lahaus/domain/usecases/users/mocks"
	"strings"
	"testing"
)

type SignUpSuite struct {
	suite.Suite
	mockCtrl      *gomock.Controller
	database      *mocks.MockStorageManager
	signUpUseCase *users.SignUpUserUseCase
}

func TestSignUpSuite(t *testing.T) {
	suite.Run(t, new(SignUpSuite))
}

func (suite *SignUpSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.database = mocks.NewMockStorageManager(suite.mockCtrl)
	conf := &config.Security{
		Secret:                 "s3cr3t",
		TokenDurationInMinutes: 1,
		Issuer:                 "lahaus",
	}
	suite.signUpUseCase = users.NewSignUpUserUseCase(conf, suite.database)
}

func (suite *SignUpSuite) TearDownSuite() {
	suite.mockCtrl.Finish()
}

func (suite *SignUpSuite) TestSignUpUseCase_ExecuteSuccess() {
	suite.database.EXPECT().GetUser(gomock.Any()).Return(&model.User{
		ID:       1,
		Email:    "d@d.com",
		Password: "a4ayc_80_OGda4BO_1o_V0etpOqiLx1JwB5S3beHW0s=",
	}, true, nil)
	token, err := suite.signUpUseCase.Execute("d@d.com", "1")
	suite.NoError(err)
	v := strings.Split(token, ".")
	suite.Len(v, 3)
}

func (suite *SignUpSuite) TestSignUpUseCase_ExecuteError_GetUser() {
	suite.database.EXPECT().GetUser(gomock.Any()).Return(nil, false, errors.New("fail to get user from database"))
	_, err := suite.signUpUseCase.Execute("d@d.com", "1")
	suite.Error(err)
}

func (suite *SignUpSuite) TestSignUpUseCase_ExecuteError_UserNotFound() {
	suite.database.EXPECT().GetUser(gomock.Any()).Return(nil, false, nil)
	_, err := suite.signUpUseCase.Execute("d@d.com", "1")
	suite.Error(err)
}

func (suite *SignUpSuite) TestSignUpUseCase_ExecuteError_PasswordDiffer() {
	suite.database.EXPECT().GetUser(gomock.Any()).Return(&model.User{
		ID:       1,
		Email:    "d@d.com",
		Password: "a4ayc_80_OGda4BO_1o_V0etpOqiLx1JwB5S3beHW0s=",
	}, true, nil)
	_, err := suite.signUpUseCase.Execute("d@d.com", "11")
	suite.Error(err)
}
