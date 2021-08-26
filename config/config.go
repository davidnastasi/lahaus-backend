package config

type Security struct {
	Secret                 string
	TokenDurationInMinutes int
	Issuer                 string
}

// SystemSettings represents the configuration of the app
type SystemSettings struct {
	Storage  *Storage
	Security *Security
	Logger   *Logger
}

// Storage represents the storage used by the app
type Storage struct {
	Database *Database
}

// Database represents the database
type Database struct {
	Host         string
	Port         uint
	User         string
	Password     string
	DatabaseName string
}

// Logger represents the logger
type Logger struct {
	Level string
}

type BetweenInt struct {
	LowerBound int
	UpperBound int
}

type BetweenFloat struct {
	LowerBound float64
	UpperBound float64
}

type PropertyTypeValidator struct {
	Bedrooms     *BetweenInt
	Bathrooms    *BetweenInt
	Area         *BetweenInt
	ParkingSpots int
}

type BundleValidator struct {
	Longitude BetweenFloat
	Latitude  BetweenFloat
	PriceIn   BetweenInt
	PriceOut  BetweenInt
}

// BusinessRules represents the business rules
type BusinessRules struct {
	HouseValidator     *PropertyTypeValidator
	ApartmentValidator *PropertyTypeValidator
	BundleValidator    *BundleValidator
}

// Config represents the configuration of system
type Config struct {
	SystemSettings *SystemSettings
	BusinessRules  *BusinessRules
}
