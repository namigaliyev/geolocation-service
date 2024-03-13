package model

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Exception struct {
	ID       string               `json:"_id"`
	Code     string               `json:"code"`
	Message  string               `json:"message"`
	HttpCode int                  `json:"httpCode"`
	Checks   *map[string][]string `json:"checks,omitempty"`
}

func (e *Exception) Error() string {
	return e.Message
}

func NewException(code, message string, httpCode int) error {
	return &Exception{
		ID:       primitive.NewObjectID().Hex(),
		Code:     code,
		Message:  message,
		HttpCode: httpCode,
		Checks:   nil,
	}
}

func NewInvalidIPAddressException(ipAddress string) error {
	return NewException(ExceptionIllegalStatusCode, InvalidIPAddressMessage+" : "+ipAddress, http.StatusBadRequest)
}

func NewGeoLocationNotFoundException(ipAddress string) error {
	return NewException(LocationNotFoundStatusCode, LocationNotFoundMessage+" : "+ipAddress, http.StatusNotFound)
}

const (
	InvalidIPAddressMessage = "Invalid IP address"
	LocationNotFoundMessage = "Location not found for given ip address"

	ExceptionIllegalStatusCode = "exception.service.illegal-argument"
	LocationNotFoundStatusCode = "exception.service.location.not-found"
)
