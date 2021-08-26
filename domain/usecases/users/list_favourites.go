package users

import (
	"lahaus/domain/model"
	"math"
)

type ListFavouritesUseCase struct {
	database StorageManager
}

type FavouritesSearchParams struct {
	Page     int64
	PageSize int64
	UserID   int64
}

func NewListFavouriteUseCase(database StorageManager) *ListFavouritesUseCase {
	return &ListFavouritesUseCase{
		database: database,
	}
}

func (uc *ListFavouritesUseCase) Execute(search FavouritesSearchParams) (*model.PropertiesPaging, error) {
	results, err := uc.database.ListFavourites(search)
	if err != nil {
		return nil, err
	}
	results.TotalPages = int64(math.Ceil(float64(results.Total) / float64(results.PageSize)))
	return results, nil
}
