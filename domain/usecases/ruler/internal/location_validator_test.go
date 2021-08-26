package internal

import (
	"lahaus/domain/model"
	"testing"
)

func TestLocationValidator_IsValidLocation(t *testing.T) {

	ruler := NewLocationRuler()
	tests := []struct {
		property  model.Property
		wantError bool
	}{
		{model.Property{Location: model.Location{Longitude: -180.0000001, Latitude: 0}}, true},
		{model.Property{Location: model.Location{Longitude: 180.0000001, Latitude: 0}}, true},
		{model.Property{Location: model.Location{Longitude: 0, Latitude: -90.0000001}}, true},
		{model.Property{Location: model.Location{Longitude: 0, Latitude: 90.0000001}}, true},
		{model.Property{Location: model.Location{Longitude: -180, Latitude: 0}}, false},
		{model.Property{Location: model.Location{Longitude: 180, Latitude: 0}}, false},
		{model.Property{Location: model.Location{Longitude: 0, Latitude: -90}}, false},
		{model.Property{Location: model.Location{Longitude: 0, Latitude: 90}}, false},
		{model.Property{Location: model.Location{Longitude: -179.9999999, Latitude: 0}}, false},
		{model.Property{Location: model.Location{Longitude: 179.9999999, Latitude: 0}}, false},
		{model.Property{Location: model.Location{Longitude: 0, Latitude: -89.9999999}}, false},
		{model.Property{Location: model.Location{Longitude: 0, Latitude: 89.9999999}}, false},
		{model.Property{Location: model.Location{Longitude: 0, Latitude: 0}}, false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := ruler(&tt.property); (got != nil) != tt.wantError {
				t.Errorf("ruler() = %v, want %v", got, tt.wantError)
			}
		})
	}
}
