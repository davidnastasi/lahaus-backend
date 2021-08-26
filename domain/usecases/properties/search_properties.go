package properties

import (
	"lahaus/domain/model"
	"math"
)

type SearchPropertyUseCase struct {
	database StorageManager
}

type BBoxSearchParams struct {
	MinLongitude float64
	MinLatitude  float64
	MaxLongitude float64
	MaxLatitude  float64
}

type PropertySearchParams struct {
	Status   string
	Bbox     *BBoxSearchParams
	Page     int64
	PageSize int64
}

func NewSearchPropertyUseCase(database StorageManager) *SearchPropertyUseCase {
	return &SearchPropertyUseCase{
		database: database,
	}
}

func (uc *SearchPropertyUseCase) Execute(search PropertySearchParams) (*model.PropertiesPaging, error) {
	propertiesPaging, err := uc.database.FilterProperties(search)
	if err != nil {
		return nil, err
	}
	propertiesPaging.TotalPages = int64(math.Ceil(float64(propertiesPaging.Total) / float64(propertiesPaging.PageSize)))
	return propertiesPaging, nil
}
