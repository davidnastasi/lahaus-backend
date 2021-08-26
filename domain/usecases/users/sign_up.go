package users

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"lahaus/config"
	"lahaus/domain/model"
	"time"
)

type SignUpUserUseCase struct {
	database StorageManager
	config   *config.Security
}

func NewSignUpUserUseCase(config *config.Security, database StorageManager) *SignUpUserUseCase {
	return &SignUpUserUseCase{
		database: database,
		config:   config,
	}
}

type UserTokenClaims struct {
	Email  string `json:"email"`
	UserID int64  `json:"userId"`
	jwt.StandardClaims
}

type Token string

func (s *SignUpUserUseCase) Execute(email, password string) (string, error) {
	user, found, err := s.database.GetUser(email)
	if err != nil {
		return "", err
	}
	if !found {
		return "", model.NewUnauthorizedError(errors.New(""))
	}

	passwordEncrypt := sha256.Sum256([]byte(password))
	passwordEncryptAsString := base64.URLEncoding.EncodeToString(passwordEncrypt[:])
	if passwordEncryptAsString != user.Password {
		return "", model.NewUnauthorizedError(errors.New("invalid credentials"))
	}
	expireTime := time.Now().Add(time.Duration(s.config.TokenDurationInMinutes) * time.Minute).Unix()
	claims := UserTokenClaims{
		email,
		user.ID,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    s.config.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", model.NewUnauthorizedError(err)
	}
	return tokenSigned, nil

}
