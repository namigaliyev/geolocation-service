package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"geolocation-service/internal/model"
)

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func HandleError(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		err := h(w, r)
		if err != nil {
			var e *model.Exception
			switch {
			case errors.As(err, &e):
				w.WriteHeader(e.HttpCode)
				json.NewEncoder(w).Encode(e)
			default:
				w.WriteHeader(http.StatusInternalServerError)
				restError := &model.Exception{
					ID:       primitive.NewObjectID().Hex(),
					Code:     http.StatusText(http.StatusInternalServerError),
					Message:  http.StatusText(http.StatusInternalServerError),
					HttpCode: http.StatusInternalServerError,
					Checks:   nil,
				}
				json.NewEncoder(w).Encode(restError)
			}
		}
	}
}
