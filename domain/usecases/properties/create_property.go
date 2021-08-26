package properties

import (
	"lahaus/domain/model"
)

//go:generate mockgen -destination=./mocks/mock_property.go -package=mocks -source=./create_property.go

type StorageManager interface {
	SaveProperty(property *model.Property) (*model.Property, error)
	UpdateProperty(property *model.Property) (*model.Property, error)
	FilterProperties(search PropertySearchParams) (*model.PropertiesPaging, error)
}

type PropertyRuler interface {
	Execute(property *model.Property)
}

type CreatePropertyUseCase struct {
	database      StorageManager
	propertyRuler PropertyRuler
}

func NewCreatePropertyUseCase(database StorageManager, propertyRuler PropertyRuler) *CreatePropertyUseCase {
	return &CreatePropertyUseCase{
		database:      database,
		propertyRuler: propertyRuler,
	}
}

func (uc CreatePropertyUseCase) Execute(property *model.Property) (*model.Property, error) {
	uc.propertyRuler.Execute(property)
	propertyStored, err := uc.database.SaveProperty(property)
	if err != nil {
		return nil, err
	}

	return propertyStored, nil

}
