package internal

import (
	"lahaus/config"
	"lahaus/domain/model"
	"testing"
)

func TestPriceValidator_IsValidPrice(t *testing.T) {

	const million = 1000 * 1000
	ruler := NewPriceRuler(&config.Config{
		BusinessRules: &config.BusinessRules{
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
	tests := []struct {
		name      string
		property  model.Property
		wantError bool
	}{
		{"property is inside and price is valid", model.Property{Location: model.Location{Longitude: -99.1, Latitude: 19.3}, Pricing: model.Pricing{SalePrice: 2 * million}}, false},
		{"property is inside and price is invalid", model.Property{Location: model.Location{Longitude: -99.1, Latitude: 19.3}, Pricing: model.Pricing{SalePrice: 20 * million}}, true},
		{"property is outside and price is valid", model.Property{Location: model.Location{Longitude: -99.5, Latitude: 19.3}, Pricing: model.Pricing{SalePrice: 200 * million}}, false},
		{"property is outside and price is invalid", model.Property{Location: model.Location{Longitude: -99.5, Latitude: 19.3}, Pricing: model.Pricing{SalePrice: 2 * million}}, true},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := ruler(&tt.property); (got != nil) != tt.wantError {
				t.Errorf("ruler() = %v, want %v", got, tt.wantError)
			}
		})
	}
}
