package internal

import (
	"errors"
	"lahaus/domain/model"
)

const minLongitude = -180.0000000
const maxLongitude = 180.0000000
const minLatitude = -90.0000000
const maxLatitude = 90.0000000

type BetweenFloat struct {
	LowerBound float64
	UpperBound float64
}

type LocationValidator struct {
	Longitude BetweenFloat
	Latitude  BetweenFloat
}

func NewLocationRuler() PropertyRulerFunc {
	validator := LocationValidator{
		Longitude: BetweenFloat{
			LowerBound: minLongitude,
			UpperBound: maxLongitude,
		},
		Latitude: BetweenFloat{
			LowerBound: minLatitude,
			UpperBound: maxLatitude,
		},
	}

	return validator.IsValidLocation()

}

func (lv *LocationValidator) IsValidLocation() PropertyRulerFunc {
	return func(property *model.Property) error {
		isValidLatitude := property.Location.Latitude >= lv.Latitude.LowerBound && property.Location.Latitude <= lv.Latitude.UpperBound
		isValidLongitude := property.Location.Longitude >= lv.Longitude.LowerBound && property.Location.Longitude <= lv.Longitude.UpperBound
		if isValidLatitude && isValidLongitude {
			return nil
		}
		return errors.New("location is not valid")
	}
}
