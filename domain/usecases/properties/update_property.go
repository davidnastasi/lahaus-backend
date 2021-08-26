package properties

import (
	"lahaus/domain/model"
)

type UpdatePropertyUseCase struct {
	database      StorageManager
	propertyRuler PropertyRuler
}

func NewUpdatePropertyUseCase(database StorageManager, propertyRuler PropertyRuler) *UpdatePropertyUseCase {
	return &UpdatePropertyUseCase{
		database:      database,
		propertyRuler: propertyRuler,
	}
}

func (uc *UpdatePropertyUseCase) Execute(property *model.Property) (*model.Property, error) {
	uc.propertyRuler.Execute(property)
	propertyStored, err := uc.database.UpdateProperty(property)
	if err != nil {
		return nil, err
	}
	return propertyStored, nil
}
