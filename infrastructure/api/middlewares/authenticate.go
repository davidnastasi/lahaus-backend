package middlewares

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"lahaus/config"
	"lahaus/domain/usecases/users"
	"net/http"
	"strings"
)

type AuthenticationMiddleware struct {
	config *config.Security
}

func NewAuthenticationMiddleware(config *config.Security) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		config: config,
	}
}

func (am *AuthenticationMiddleware) Execute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value, ok := r.Header["Authorization"]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		receivedToken := value[0]
		receivedToken = strings.ReplaceAll(receivedToken, "Bearer ", "")
		claims := &users.UserTokenClaims{}
		token, err := jwt.ParseWithClaims(receivedToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(am.config.Secret), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		values := map[string]interface{}{
			"email":  claims.Email,
			"userId": claims.UserID,
		}
		ctx := context.WithValue(r.Context(), "user", values)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
