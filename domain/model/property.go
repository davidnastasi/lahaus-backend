package model

import "time"

type PropertyStatus string

const (
	ACTIVE   PropertyStatus = "ACTIVE"
	INACTIVE PropertyStatus = "INACTIVE"
	INVALID  PropertyStatus = "INVALID"
)

type PropertyType string

const (
	HOUSE     PropertyType = "HOUSE"
	APARTMENT PropertyType = "APARTMENT"
)

type Photos []string

type ParkingSpots *int
type AdministrativeFee *int
type Description *string

type Property struct {
	ID           int64          `json:"id"`
	Title        string         `json:"title"`
	Description  Description    `json:"description,omitempty"`
	Location     Location       `json:"location"`
	Pricing      Pricing        `json:"pricing"`
	PropertyType PropertyType   `json:"propertyType"`
	Bedrooms     int            `json:"bedrooms"`
	Bathrooms    int            `json:"bathrooms"`
	ParkingSpots ParkingSpots   `json:"parkingSpots,omitempty"`
	Area         int            `json:"area"`
	Photos       Photos         `json:"photos,omitempty"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	Status       PropertyStatus `json:"status"`
}
type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
type Pricing struct {
	SalePrice         int               `json:"salePrice"`
	AdministrativeFee AdministrativeFee `json:"administrativeFee,omitempty"`
}

type PropertiesPaging struct {
	Page       int64       `json:"page"`
	PageSize   int64       `json:"pageSize"`
	TotalPages int64       `json:"totalPages"`
	Total      int64       `json:"total"`
	Data       []*Property `json:"data"`
}
