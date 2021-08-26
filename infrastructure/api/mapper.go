package api

import (
	"fmt"
	"lahaus/domain/model"
	"net/mail"
	"strings"
)

const HOUSE = "HOUSE"
const APARTMENT = "APARTMENT"

func mapPropertyRequestToProperty(request propertyRequest) (*model.Property, error) {
	propertyType, err := mapStringToPropertyType(request.PropertyType)
	if err != nil {
		return nil, err
	}
	return &model.Property{
		Title:       request.Title,
		Description: request.Description,
		Location: model.Location{
			Longitude: request.Location.Longitude,
			Latitude:  request.Location.Latitude,
		},
		Pricing: model.Pricing{
			SalePrice:         request.Pricing.SalePrice,
			AdministrativeFee: request.Pricing.AdministrativeFee,
		},
		PropertyType: propertyType,
		Bedrooms:     request.Bedrooms,
		Bathrooms:    request.Bathrooms,
		ParkingSpots: request.ParkingSpots,
		Area:         request.Area,
		Photos:       request.Photos,
	}, nil
}

func mapStringToPropertyType(propertyTypeAsString string) (model.PropertyType, error) {
	input := strings.ToUpper(propertyTypeAsString)
	switch input {
	case HOUSE:
		return model.HOUSE, nil
	case APARTMENT:
		return model.APARTMENT, nil
	default:
		return model.HOUSE, fmt.Errorf("property type not recognized [%s]", input)
	}
}

func mapCreateUserRequestToUser(request createUserRequest) (*model.User, error) {
	_, err := mail.ParseAddress(request.Email)
	if err != nil {
		return nil, err
	}
	return &model.User{
		Email:    request.Email,
		Password: request.Email,
	}, nil
}
