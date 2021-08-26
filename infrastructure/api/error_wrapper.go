package api

import (
	"encoding/json"
	"go.uber.org/zap"
	"lahaus/domain/model"
	"lahaus/logger"
	"net/http"
)

func wrapError(w http.ResponseWriter, err error, code int) {
	switch err.(type) {
	case *model.DomainError:
		responseWriter(w, err, http.StatusBadRequest)
		return
	case *model.EntityNotFoundError:
		responseWriter(w, err, http.StatusNotFound)
		return
	case *model.InternalServerError:
		responseWriter(w, err, http.StatusInternalServerError)
		return
	case *model.UnauthorizedError:
		responseWriter(w, err, http.StatusUnauthorized)
		return
	}

	switch code {
	case http.StatusUnauthorized:
		responseWriter(w, model.NewUnauthorizedError(err), code)
	case http.StatusBadRequest:
		responseWriter(w, model.NewDomainError(err), code)
	case http.StatusInternalServerError:
		responseWriter(w, model.NewInternalServerError(err), code)
	}
}

func responseWriter(w http.ResponseWriter, err error, code int) {
	response, err := json.Marshal(err)
	if err != nil {
		logger.GetInstance().Error("error marshalling response error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		logger.GetInstance().Error("error writing response error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
