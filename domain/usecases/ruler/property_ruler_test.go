package ruler

import (
	"lahaus/config"
	"lahaus/domain/model"
	"testing"
)

func TestPropertyRuler_Execute(t *testing.T) {

	const million = 1000 * 1000
	ruler := NewPropertyRulerUseCase(&config.Config{
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

	value1 := 1
	tests := []struct {
		name       string
		property   model.Property
		wantStatus model.PropertyStatus
	}{

		{"property house is invalid", model.Property{PropertyType: model.HOUSE, Bedrooms: 0, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value1), Area: 400}, model.INVALID},
		{"property apartment is invalid", model.Property{PropertyType: model.APARTMENT, Bedrooms: 0, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value1), Area: 400}, model.INVALID},
		{"property has valid features, is inside and price is valid", model.Property{
			PropertyType: model.HOUSE, Bedrooms: 2, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value1), Area: 400,
			Location: model.Location{Longitude: -99.1, Latitude: 19.3}, Pricing: model.Pricing{SalePrice: 2 * million}}, model.ACTIVE},
		{"property has valid features, is inside and price is invalid", model.Property{
			PropertyType: model.HOUSE, Bedrooms: 2, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value1), Area: 400,
			Location: model.Location{Longitude: -99.1, Latitude: 19.3}, Pricing: model.Pricing{SalePrice: 20 * million}}, model.INVALID},
		{"property has valid features, is outside and price is valid", model.Property{
			PropertyType: model.HOUSE, Bedrooms: 2, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value1), Area: 400,
			Location: model.Location{Longitude: -99.5, Latitude: 19.3}, Pricing: model.Pricing{SalePrice: 200 * million}}, model.INACTIVE},
		{"property has valid features, is valid outside and price is invalid", model.Property{
			PropertyType: model.HOUSE, Bedrooms: 2, Bathrooms: 2, ParkingSpots: model.ParkingSpots(&value1), Area: 400,
			Location: model.Location{Longitude: -99.5, Latitude: 19.3}, Pricing: model.Pricing{SalePrice: 2 * million}}, model.INVALID},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			ruler.Execute(&tt.property)
			if tt.property.Status != tt.wantStatus {
				t.Errorf("ruler() = %v, want %v", tt.property.Status, tt.wantStatus)
			}
		})
	}

}
