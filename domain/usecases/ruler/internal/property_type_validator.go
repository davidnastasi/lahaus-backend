package internal

import (
	"fmt"
	"lahaus/config"
	"lahaus/domain/model"
)

type PropertyTypeRuler map[model.PropertyType]PropertyRulerFuncs

type BetweenInt struct {
	LowerBound int
	UpperBound int
}

type PropertyTypeRules struct {
	Bedrooms     BetweenInt
	Bathrooms    BetweenInt
	Area         BetweenInt
	ParkingSpots int
}

func NewPropertyTypeRuler(config *config.Config) PropertyRulerFunc {
	houseValidator := PropertyTypeRules{
		Bedrooms: BetweenInt{
			LowerBound: config.BusinessRules.HouseValidator.Bedrooms.LowerBound,
			UpperBound: config.BusinessRules.HouseValidator.Bedrooms.UpperBound,
		},
		Bathrooms: BetweenInt{
			LowerBound: config.BusinessRules.HouseValidator.Bathrooms.LowerBound,
			UpperBound: config.BusinessRules.HouseValidator.Bathrooms.UpperBound,
		},
		Area: BetweenInt{
			LowerBound: config.BusinessRules.HouseValidator.Area.LowerBound,
			UpperBound: config.BusinessRules.HouseValidator.Area.UpperBound,
		},
		ParkingSpots: config.BusinessRules.HouseValidator.ParkingSpots,
	}
	apartmentValidator := PropertyTypeRules{
		Bedrooms: BetweenInt{
			LowerBound: config.BusinessRules.ApartmentValidator.Bedrooms.LowerBound,
			UpperBound: config.BusinessRules.ApartmentValidator.Bedrooms.UpperBound,
		},
		Bathrooms: BetweenInt{
			LowerBound: config.BusinessRules.ApartmentValidator.Bathrooms.LowerBound,
			UpperBound: config.BusinessRules.ApartmentValidator.Bathrooms.UpperBound,
		},
		Area: BetweenInt{
			LowerBound: config.BusinessRules.ApartmentValidator.Area.LowerBound,
			UpperBound: config.BusinessRules.ApartmentValidator.Area.UpperBound,
		},
		ParkingSpots: config.BusinessRules.ApartmentValidator.ParkingSpots,
	}

	houseRuleFunc := PropertyRulerFuncs{
		houseValidator.IsValidBedrooms(),
		houseValidator.IsValidBathrooms(),
		houseValidator.IsValidArea(),
		houseValidator.IsValidParkingSpot(),
	}

	apartmentRuleFunc := PropertyRulerFuncs{
		apartmentValidator.IsValidBedrooms(),
		apartmentValidator.IsValidBathrooms(),
		apartmentValidator.IsValidArea(),
		apartmentValidator.IsValidParkingSpot(),
	}

	ruler := PropertyTypeRuler{
		model.HOUSE:     houseRuleFunc,
		model.APARTMENT: apartmentRuleFunc,
	}

	return func(property *model.Property) error {
		rules := ruler[property.PropertyType]
		for _, fn := range rules {
			if err := fn(property); err != nil {
				return err
			}
		}
		return nil
	}

}

func (ptv PropertyTypeRules) IsValidBedrooms() PropertyRulerFunc {
	return func(property *model.Property) error {
		if property.Bedrooms < ptv.Bedrooms.LowerBound || property.Bedrooms > ptv.Bedrooms.UpperBound {
			return fmt.Errorf("bedrooms must be between %v and %v", ptv.Bedrooms.LowerBound, ptv.Bedrooms.UpperBound)
		}
		return nil
	}
}

func (ptv PropertyTypeRules) IsValidBathrooms() PropertyRulerFunc {
	return func(property *model.Property) error {
		if property.Bathrooms < ptv.Bathrooms.LowerBound || property.Bathrooms > ptv.Bathrooms.UpperBound {
			return fmt.Errorf("bathrooms must be between %v and %v", ptv.Bathrooms.LowerBound, ptv.Bathrooms.UpperBound)
		}
		return nil
	}
}

func (ptv PropertyTypeRules) IsValidArea() PropertyRulerFunc {
	return func(property *model.Property) error {
		if property.Area < ptv.Area.LowerBound || property.Area > ptv.Area.UpperBound {
			return fmt.Errorf("area must be between %v and %v", ptv.Area.LowerBound, ptv.Area.UpperBound)
		}
		return nil
	}
}

func (ptv PropertyTypeRules) IsValidParkingSpot() PropertyRulerFunc {
	return func(property *model.Property) error {
		if property.ParkingSpots != nil && (*property.ParkingSpots < ptv.ParkingSpots) {
			return fmt.Errorf("parkingspots must be greater than %v", ptv.ParkingSpots)
		}
		return nil
	}
}
