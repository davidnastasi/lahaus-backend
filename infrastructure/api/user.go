package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"lahaus/domain/model"
	"lahaus/domain/usecases/users"
	"lahaus/logger"
	"net/http"
	"net/url"
	"strconv"
)

type createUserRequest struct {
	Email string `json:"email"`
}

type signUpUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signUpUserResponse struct {
	Token string `json:"token"`
}

//go:generate mockgen -destination=./mocks/mock_user.go -package=mocks -source=./user.go

type SignInUserExecutor interface {
	Execute(user *model.User) error
}

type SignUpUserExecutor interface {
	Execute(email, password string) (string, error)
}

type AddFavouriteExecutor interface {
	Execute(userID int64, property int64) error
}

type ListFavouritesExecutor interface {
	Execute(search users.FavouritesSearchParams) (*model.PropertiesPaging, error)
}

type UserHandler struct {
	createUserExecutor     SignInUserExecutor
	signUpUserExecutor     SignUpUserExecutor
	addFavouriteExecutor   AddFavouriteExecutor
	listFavouritesExecutor ListFavouritesExecutor
}

func NewUserHandler(createUserExecutor SignInUserExecutor, signUpUserExecutor SignUpUserExecutor, addFavouriteExecutor AddFavouriteExecutor, listFavouritesExecutor ListFavouritesExecutor) *UserHandler {
	return &UserHandler{
		createUserExecutor:     createUserExecutor,
		signUpUserExecutor:     signUpUserExecutor,
		addFavouriteExecutor:   addFavouriteExecutor,
		listFavouritesExecutor: listFavouritesExecutor,
	}
}

// SignInUser  handler the request
func (handler *UserHandler) SignInUser(w http.ResponseWriter, r *http.Request) {
	var request createUserRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logger.GetInstance().Error("json decode error", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusBadRequest)
		return
	}

	user, err := mapCreateUserRequestToUser(request)
	if err != nil {
		logger.GetInstance().Error("error mapping to property", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusBadRequest)
		return
	}

	err = handler.createUserExecutor.Execute(user)
	if err != nil {
		logger.GetInstance().Error("error validating input", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

// SignUpUser  handler the request
func (handler *UserHandler) SignUpUser(w http.ResponseWriter, r *http.Request) {
	var request signUpUserRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logger.GetInstance().Error("json decode error", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}

	token, err := handler.signUpUserExecutor.Execute(request.Email, request.Password)
	if err != nil {
		logger.GetInstance().Error("error in login", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusUnauthorized)
		return
	}

	response := signUpUserResponse{Token: token}
	responseJson, err := json.Marshal(response)
	if err != nil {
		logger.GetInstance().Error("error in marshalling login response", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(responseJson)
	if err != nil {
		logger.GetInstance().Error("error in write login response", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type AddFavouriteToUserRequest struct {
	PropertyID int64 `json:"propertyId"`
}

// AddFavourite  handler the request
func (handler *UserHandler) AddFavourite(w http.ResponseWriter, r *http.Request) {
	var request AddFavouriteToUserRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logger.GetInstance().Error("json decode error", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusBadRequest)
		return
	}

	ctxUser := r.Context().Value("user")
	values := ctxUser.(map[string]interface{})
	id := values["userId"].(int64)

	err = handler.addFavouriteExecutor.Execute(id, request.PropertyID)
	if err != nil {
		logger.GetInstance().Error("error adding favourite", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

// ListFavourites  handler the request
func (handler *UserHandler) ListFavourites(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(map[string]interface{})["userId"].(int64)
	searchParams, err := mapToFavouriteSearchParams(r.URL.Query())
	if err != nil {
		logger.GetInstance().Error("error listing favourites", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusBadRequest)
		return
	}
	searchParams.UserID = userId

	results, err := handler.listFavouritesExecutor.Execute(searchParams)
	if err != nil {
		logger.GetInstance().Error("error listing favourites", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}

	responseJson, err := json.Marshal(results)
	if err != nil {
		logger.GetInstance().Error("error in marshalling favourites response", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseJson)
	if err != nil {
		logger.GetInstance().Error("error in write login response", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func mapToFavouriteSearchParams(query url.Values) (users.FavouritesSearchParams, error) {
	searchParams := users.FavouritesSearchParams{}

	page := query.Get("page")
	if page != "" {
		pageValues, err := strconv.ParseInt(page, 10, 64)
		if err != nil {
			return searchParams, err
		}
		searchParams.Page = pageValues
	} else {
		searchParams.Page = 1
	}

	pageSize := query.Get("pageSize")
	if pageSize != "" {
		pageSizeValue, err := strconv.ParseInt(pageSize, 10, 64)
		if err != nil {
			return searchParams, err
		}
		if pageSizeValue < 10 || pageSizeValue > 20 {
			return searchParams, err
		}

		searchParams.PageSize = pageSizeValue
	} else {
		searchParams.PageSize = 10
	}

	return searchParams, nil

}
