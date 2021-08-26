package internal

import (
	"github.com/stretchr/testify/require"
	"lahaus/config"
	"lahaus/domain/model"
	"testing"
)

func TestPropertyTypeRuler(t *testing.T) {
	ruler := NewPropertyTypeRuler(&config.Config{
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
		},
	})

	value1 := 1
	value2 := -1
	value3 := 0

	validHouseProperty := &model.Property{
		PropertyType: model.HOUSE,
		Bedrooms:     3,
		Bathrooms:    2,
		ParkingSpots: model.ParkingSpots(&value1),
		Area:         400,
	}

	require.NoError(t, ruler(validHouseProperty))

	invalidHouseProperties := []*model.Property{
		{PropertyType: model.HOUSE, Bedrooms: 0, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value1), Area: 400},
		{PropertyType: model.HOUSE, Bedrooms: 100, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value1), Area: 400},
		{PropertyType: model.HOUSE, Bedrooms: 3, Bathrooms: 0, ParkingSpots: model.ParkingSpots(&value1), Area: 400},
		{PropertyType: model.HOUSE, Bedrooms: 3, Bathrooms: 100, ParkingSpots: model.ParkingSpots(&value1), Area: 400},
		{PropertyType: model.HOUSE, Bedrooms: 3, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value1), Area: 40},
		{PropertyType: model.HOUSE, Bedrooms: 3, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value1), Area: 5000},
		{PropertyType: model.HOUSE, Bedrooms: 3, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value2), Area: 400},
	}

	for _, invalidHouseProperty := range invalidHouseProperties {
		require.Error(t, ruler(invalidHouseProperty))
	}

	validApartmentProperty := &model.Property{
		PropertyType: model.APARTMENT,
		Bedrooms:     3,
		Bathrooms:    2,
		ParkingSpots: model.ParkingSpots(&value1),
		Area:         400,
	}

	require.NoError(t, ruler(validApartmentProperty))

	invalidApartmentProperties := []*model.Property{
		{PropertyType: model.APARTMENT, Bedrooms: 0, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value1), Area: 400},
		{PropertyType: model.APARTMENT, Bedrooms: 100, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value1), Area: 400},
		{PropertyType: model.APARTMENT, Bedrooms: 3, Bathrooms: 0, ParkingSpots: model.ParkingSpots(&value1), Area: 400},
		{PropertyType: model.APARTMENT, Bedrooms: 3, Bathrooms: 100, ParkingSpots: model.ParkingSpots(&value1), Area: 400},
		{PropertyType: model.APARTMENT, Bedrooms: 3, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value1), Area: 30},
		{PropertyType: model.APARTMENT, Bedrooms: 3, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value1), Area: 5000},
		{PropertyType: model.APARTMENT, Bedrooms: 3, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value3), Area: 400},
	}

	for _, invalidApartmentProperty := range invalidApartmentProperties {
		require.Error(t, ruler(invalidApartmentProperty))
	}

}
