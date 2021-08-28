package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"lahaus/domain/model"
	"lahaus/domain/usecases/properties"
	"lahaus/logger"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//go:generate mockgen -destination=./mocks/mock_property.go -package=mocks -source=./property.go

// PropertyExecutor ...
type PropertyExecutor interface {
	Execute(property *model.Property) (*model.Property, error)
}

// SearchPropertyExecutor ...
type SearchPropertyExecutor interface {
	Execute(search properties.PropertySearchParams) (*model.PropertiesPaging, error)
}

const minLongitudeValue = -180.0000000
const maxLongitudeValue = 180.0000000
const minLatitudeValue = -90.0000000
const maxLatitudeValue = 90.0000000

// PropertyHandler struct
type PropertyHandler struct {
	createPropertyExecutor PropertyExecutor
	updatePropertyExecutor PropertyExecutor
	searchExecutor         SearchPropertyExecutor
}

// NewPropertyHandler creates a new PropertyHandler
func NewPropertyHandler(createExecutor, updateExecutor PropertyExecutor, filterExecutor SearchPropertyExecutor) *PropertyHandler {
	return &PropertyHandler{
		createPropertyExecutor: createExecutor,
		updatePropertyExecutor: updateExecutor,
		searchExecutor:         filterExecutor,
	}
}

type propertyRequest struct {
	Title        *string  `json:"title"`
	Description  *string  `json:"description"`
	Location     location `json:"location"`
	Pricing      pricing  `json:"pricing"`
	PropertyType string   `json:"propertyType"`
	Bedrooms     *int     `json:"bedrooms"`
	Bathrooms    *int     `json:"bathrooms"`
	ParkingSpots *int     `json:"parkingSpots"`
	Area         *int     `json:"area"`
	Photos       []string `json:"photos"`
}

type location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type pricing struct {
	SalePrice         *int `json:"salePrice"`
	AdministrativeFee *int `json:"administrativeFee"`
}

// CreateProperty property handler the request
func (handler *PropertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	var request propertyRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logger.GetInstance().Error("json decode error", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusBadRequest)
		return
	}

	property, err := mapPropertyRequestToProperty(request)
	if err != nil {
		logger.GetInstance().Error("error mapping to property", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusBadRequest)
		return
	}

	property, err = handler.createPropertyExecutor.Execute(property)
	if err != nil {
		logger.GetInstance().Error("error creating property", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(property)
	if err != nil {
		logger.GetInstance().Error("error marshalling property", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		logger.GetInstance().Error("error writing response", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

// UpdateProperty property handler the request
func (handler *PropertyHandler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	var request propertyRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logger.GetInstance().Error("json decode error", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusBadRequest)
		return
	}

	idValue := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idValue, 10, 64)
	if err != nil {
		logger.GetInstance().Error("error in parsing id ", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusBadRequest)
		return
	}

	property, err := mapPropertyRequestToProperty(request)
	if err != nil {
		logger.GetInstance().Error("error mapping to property", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusBadRequest)
		return
	}

	property.ID = id
	property, err = handler.updatePropertyExecutor.Execute(property)
	if err != nil {
		logger.GetInstance().Error("error updating property", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(property)
	if err != nil {
		logger.GetInstance().Error("error marshalling property", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		logger.GetInstance().Error("error writing response", zap.Error(err), zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

// SearchProperties property handler the request
func (handler *PropertyHandler) SearchProperties(w http.ResponseWriter, r *http.Request) {
	searchParams, err := mapToPropertySearchParams(r.URL.Query())
	if err != nil {
		logger.GetInstance().Error("error validating input", zap.Error(err),
			zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusBadRequest)
		return
	}

	results, err := handler.searchExecutor.Execute(searchParams)
	if err != nil {
		logger.GetInstance().Error("error getting results", zap.Error(err),
			zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(results)
	if err != nil {
		logger.GetInstance().Error("error marshalling results ", zap.Error(err),
			zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		logger.GetInstance().Error("error writing response", zap.Error(err),
			zap.String(middleware.RequestIDHeader, r.Context().Value(middleware.RequestIDKey).(string)))
		wrapError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func mapToPropertySearchParams(query url.Values) (properties.PropertySearchParams, error) {
	searchParams := properties.PropertySearchParams{}
	status := query.Get("status")
	if status == "" {
		searchParams.Status = "ALL"
	} else {
		if status != "ALL" && status != "ACTIVE" && status != "INACTIVE" && status != "INVALID" {
			return searchParams, fmt.Errorf("invalid status [%v]", status)
		}
		searchParams.Status = status
	}
	bbox := query.Get("bbox")
	if bbox != "" {
		bboxNormalized := strings.ReplaceAll(bbox, " ", "")
		bboxValues := strings.Split(bboxNormalized, ",")
		if len(bboxValues) != 4 {
			return searchParams, fmt.Errorf("invalid bbox format [%v]", bbox)
		}

		minLongitude, err := strconv.ParseFloat(bboxValues[0], 64)
		if err != nil {
			return searchParams, err
		}

		minLatitude, err := strconv.ParseFloat(bboxValues[1], 64)
		if err != nil {
			return searchParams, err
		}

		maxLongitude, err := strconv.ParseFloat(bboxValues[2], 64)
		if err != nil {
			return searchParams, err
		}

		maxLatitude, err := strconv.ParseFloat(bboxValues[3], 64)
		if err != nil {
			return searchParams, err
		}

		isValidLongitude := minLongitude >= minLongitudeValue && maxLongitude <= maxLongitudeValue
		isValidLatitude := minLatitude >= minLatitudeValue && maxLatitude <= maxLatitudeValue
		if !isValidLatitude || !isValidLongitude {
			return searchParams, errors.New("location is not valid")
		}

		searchParams.Bbox = &properties.BBoxSearchParams{
			MinLongitude: minLongitude,
			MinLatitude:  minLatitude,
			MaxLongitude: maxLongitude,
			MaxLatitude:  maxLatitude,
		}
	}
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
			return searchParams, errors.New("page size should be between 10 and 20")
		}

		searchParams.PageSize = pageSizeValue
	} else {
		searchParams.PageSize = 10
	}

	return searchParams, nil

}
