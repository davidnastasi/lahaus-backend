package properties_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"lahaus/config"
	"lahaus/domain/model"
	"lahaus/domain/usecases/properties"
	"lahaus/domain/usecases/properties/mocks"
	"lahaus/domain/usecases/ruler"
	"testing"
)

const million = 1000000

type CreatePropertySuite struct {
	suite.Suite
	mockCtrl      *gomock.Controller
	database      *mocks.MockStorageManager
	propertyRuler *ruler.PropertyRules
	createUseCase *properties.CreatePropertyUseCase
}

func TestCreatePropertySuite(t *testing.T) {
	suite.Run(t, new(CreatePropertySuite))
}

func (suite *CreatePropertySuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.database = mocks.NewMockStorageManager(suite.mockCtrl)
	suite.propertyRuler = ruler.NewPropertyRulerUseCase(&config.Config{
		BusinessRules: &config.BusinessRules{
			HouseValidator: &config.PropertyTypeValidator{
				Bedrooms: &config.BetweenInt{
					LowerBound: 1,
					UpperBound: 14,
				},
				Bathrooms: &config.BetweenInt{
					LowerBound: 1,
					UpperBound: 12,
				},
				ParkingSpots: 0,
				Area: &config.BetweenInt{
					LowerBound: 50,
					UpperBound: 3000,
				},
			},
			ApartmentValidator: &config.PropertyTypeValidator{
				Bedrooms: &config.BetweenInt{
					LowerBound: 1,
					UpperBound: 6,
				},
				Bathrooms: &config.BetweenInt{
					LowerBound: 1,
					UpperBound: 4,
				},
				ParkingSpots: 1,
				Area: &config.BetweenInt{
					LowerBound: 40,
					UpperBound: 400,
				},
			},
			BundleValidator: &config.BundleValidator{
				Longitude: config.BetweenFloat{
					LowerBound: -99.296741,
					UpperBound: -98.916339,
				},
				Latitude: config.BetweenFloat{
					LowerBound: 19.296134,
					UpperBound: 19.661237,
				},
				PriceIn: config.BetweenInt{
					LowerBound: million,
					UpperBound: 15 * million,
				},
				PriceOut: config.BetweenInt{
					LowerBound: 50 * million,
					UpperBound: 3500 * million,
				},
			},
		},
	})
	suite.createUseCase = properties.NewCreatePropertyUseCase(suite.database, suite.propertyRuler)
}

func (suite *CreatePropertySuite) TearDownSuite() {
	suite.mockCtrl.Finish()
}

func (suite *CreatePropertySuite) TestCreatePropertyUseCase_ExecuteSuccessActive() {
	property := &model.Property{
		Title:       "Casa de familia",
		Description: nil,
		Location: model.Location{
			Longitude: -99.096741,
			Latitude:  19.296135,
		},
		Pricing: model.Pricing{
			SalePrice:         3 * million,
			AdministrativeFee: nil,
		},
		PropertyType: model.HOUSE,
		Bedrooms:     1,
		Bathrooms:    1,
		ParkingSpots: nil,
		Area:         300,
		Photos:       nil,
	}

	suite.database.EXPECT().SaveProperty(property).Return(property, nil)
	propertyResult, err := suite.createUseCase.Execute(property)
	suite.NoError(err)
	suite.Equal(model.ACTIVE, propertyResult.Status)
}

func (suite *CreatePropertySuite) TestCreatePropertyUseCase_ExecuteSuccessInactive() {
	property := &model.Property{
		Title:       "Casa de familia",
		Description: nil,
		Location: model.Location{
			Longitude: -99.096741,
			Latitude:  20.296135,
		},
		Pricing: model.Pricing{
			SalePrice:         130 * million,
			AdministrativeFee: nil,
		},
		PropertyType: model.HOUSE,
		Bedrooms:     1,
		Bathrooms:    1,
		ParkingSpots: nil,
		Area:         300,
		Photos:       nil,
	}

	suite.database.EXPECT().SaveProperty(property).Return(property, nil)
	propertyResult, err := suite.createUseCase.Execute(property)
	suite.NoError(err)
	suite.Equal(model.INACTIVE, propertyResult.Status)
}

func (suite *CreatePropertySuite) TestCreatePropertyUseCase_ExecuteSuccessInvalidAmount() {
	property := &model.Property{
		Title:       "Casa de familia",
		Description: nil,
		Location: model.Location{
			Longitude: -99.096741,
			Latitude:  20.296135,
		},
		Pricing: model.Pricing{
			SalePrice:         30 * million,
			AdministrativeFee: nil,
		},
		PropertyType: model.HOUSE,
		Bedrooms:     1,
		Bathrooms:    1,
		ParkingSpots: nil,
		Area:         300,
		Photos:       nil,
	}

	suite.database.EXPECT().SaveProperty(property).Return(property, nil)
	propertyResult, err := suite.createUseCase.Execute(property)
	suite.NoError(err)
	suite.Equal(model.INVALID, propertyResult.Status)
}

func (suite *CreatePropertySuite) TestCreatePropertyUseCase_ExecuteError() {
	property := &model.Property{
		Title:       "Casa de familia",
		Description: nil,
		Location: model.Location{
			Longitude: -99.096741,
			Latitude:  220.296135,
		},
		Pricing: model.Pricing{
			SalePrice:         130 * million,
			AdministrativeFee: nil,
		},
		PropertyType: model.HOUSE,
		Bedrooms:     1,
		Bathrooms:    1,
		ParkingSpots: nil,
		Area:         300,
		Photos:       nil,
	}

	suite.database.EXPECT().SaveProperty(property).Return(nil, errors.New("fail to save in database"))
	propertyResult, err := suite.createUseCase.Execute(property)
	suite.Error(err)
	suite.Nil(propertyResult)
}
