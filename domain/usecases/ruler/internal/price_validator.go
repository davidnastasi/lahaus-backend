package internal

import (
	"fmt"
	"lahaus/config"
	"lahaus/domain/model"
)

type PriceValidator struct {
	Longitude BetweenFloat
	Latitude  BetweenFloat
	PriceIn   BetweenInt
	PriceOut  BetweenInt
}

func NewPriceRuler(config *config.Config) PropertyRulerFunc {
	validator := PriceValidator{
		PriceIn: BetweenInt{
			LowerBound: config.BusinessRules.BundleValidator.PriceIn.LowerBound,
			UpperBound: config.BusinessRules.BundleValidator.PriceIn.UpperBound,
		},
		PriceOut: BetweenInt{
			LowerBound: config.BusinessRules.BundleValidator.PriceOut.LowerBound,
			UpperBound: config.BusinessRules.BundleValidator.PriceOut.UpperBound,
		},
		Longitude: BetweenFloat{
			LowerBound: config.BusinessRules.BundleValidator.Longitude.LowerBound,
			UpperBound: config.BusinessRules.BundleValidator.Longitude.UpperBound,
		},
		Latitude: BetweenFloat{
			LowerBound: config.BusinessRules.BundleValidator.Latitude.LowerBound,
			UpperBound: config.BusinessRules.BundleValidator.Latitude.UpperBound,
		},
	}

	return validator.IsValidPrice()

}

func (lv *PriceValidator) IsValidPrice() PropertyRulerFunc {
	return func(property *model.Property) error {
		inside := lv.IsInsideBundleBox(property)

		if inside {
			property.Status = model.ACTIVE
			if property.Pricing.SalePrice < lv.PriceIn.LowerBound || property.Pricing.SalePrice > lv.PriceIn.UpperBound {
				return fmt.Errorf("price is incorrect") // TODO: return values
			}
			return nil
		}
		property.Status = model.INACTIVE
		if property.Pricing.SalePrice < lv.PriceOut.LowerBound || property.Pricing.SalePrice > lv.PriceOut.UpperBound {
			return fmt.Errorf("price is incorrect") //TODO: return values
		}
		return nil
	}
}

func (lv *PriceValidator) IsInsideBundleBox(property *model.Property) bool {
	return lv.isInsideLatitudeBox(property) && lv.isInsideLongitudeBox(property)
}

func (lv *PriceValidator) isInsideLatitudeBox(property *model.Property) bool {
	return property.Location.Latitude >= lv.Latitude.LowerBound && property.Location.Latitude <= lv.Latitude.UpperBound
}

func (lv *PriceValidator) isInsideLongitudeBox(property *model.Property) bool {
	return property.Location.Longitude >= lv.Longitude.LowerBound && property.Location.Longitude <= lv.Longitude.UpperBound
}
