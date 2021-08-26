package users

import (
	"crypto/sha256"
	"encoding/base64"
	"lahaus/domain/model"
)

//go:generate mockgen -destination=./mocks/mock_signin.go -package=mocks -source=./sign_in.go

type StorageManager interface {
	SaveUser(user *model.User) error
	GetUser(emil string) (*model.User, bool, error)
	GetProperty(id int64) (*model.Property, bool, error)
	AddFavourite(userID, propertyID int64) error
	ListFavourites(search FavouritesSearchParams) (*model.PropertiesPaging, error)
}

type SignInUserUseCase struct {
	database StorageManager
}

func NewSignInUserUseCase(database StorageManager) *SignInUserUseCase {
	return &SignInUserUseCase{
		database: database,
	}
}

func (c *SignInUserUseCase) Execute(user *model.User) error {
	passwordEncrypt := sha256.Sum256([]byte(user.Password))
	user.Password = base64.URLEncoding.EncodeToString(passwordEncrypt[:])
	return c.database.SaveUser(user)
}
