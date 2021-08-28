package api

import (
	"errors"
	"fmt"
	"lahaus/domain/model"
	"net/mail"
	"strings"
)

const HOUSE = "HOUSE"
const APARTMENT = "APARTMENT"

func mapPropertyRequestToProperty(request propertyRequest) (*model.Property, error) {
	if request.Title == nil || len(*request.Title) == 0 {
		return nil, errors.New("title field is a must")
	}
	if request.Bedrooms == nil {
		return nil, errors.New("bedrooms field is a must")
	}

	if request.Bathrooms == nil {
		return nil, errors.New("bathrooms field is a must")
	}

	if request.Pricing.SalePrice == nil {
		return nil, errors.New("salePrice field is a must")
	}
	if request.Area == nil {
		return nil, errors.New("area field is a must")
	}
	propertyType, err := mapStringToPropertyType(request.PropertyType)
	if err != nil {
		return nil, err
	}
	return &model.Property{
		Title:       *request.Title,
		Description: request.Description,
		Location: model.Location{
			Longitude: request.Location.Longitude,
			Latitude:  request.Location.Latitude,
		},
		Pricing: model.Pricing{
			SalePrice:         *request.Pricing.SalePrice,
			AdministrativeFee: request.Pricing.AdministrativeFee,
		},
		PropertyType: propertyType,
		Bedrooms:     *request.Bedrooms,
		Bathrooms:    *request.Bathrooms,
		ParkingSpots: request.ParkingSpots,
		Area:         *request.Area,
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
