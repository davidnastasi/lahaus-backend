package users

import (
	"errors"
	"lahaus/domain/model"
)

type AddFavouriteUseCase struct {
	database StorageManager
}

func NewAddFavouriteUseCase(database StorageManager) *AddFavouriteUseCase {
	return &AddFavouriteUseCase{
		database: database,
	}
}

func (uc *AddFavouriteUseCase) Execute(userID int64, propertyID int64) error {
	_, found, err := uc.database.GetProperty(propertyID)
	if err != nil {
		return err
	}
	if !found {
		return model.NewEntityNotFoundError(errors.New("property not found"))
	}
	err = uc.database.AddFavourite(userID, propertyID)
	if err != nil {
		return err
	}
	return nil
}
