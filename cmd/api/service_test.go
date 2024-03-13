package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"geolocation-service/config"
	"geolocation-service/internal/model"
	"geolocation-service/mocks"
)

var mockMongoDB = new(mocks.NoSQLDb)
var cfg = config.Config{
	MongoDB: mockMongoDB,
}
var service = &GeolocationService{cfg: cfg}

func TestGetGeolocationByIpAddress_Success(t *testing.T) {
	now := time.Now()
	locationID := primitive.NewObjectID()
	expectedGeolocation := &model.Geolocation{
		City:         "Baku",
		CountryCode:  "AZ",
		CreatedAt:    now,
		ID:           locationID,
		IpAddress:    "192.168.1.1",
		Latitude:     36.6532,
		Longitude:    86.5035,
		MysteryValue: 2421,
	}
	mockMongoDB.
		On("GetGeolocationByIpAddress", mock.Anything, "192.168.1.1").
		Return(expectedGeolocation, nil)

	geolocation, err := service.
		GetGeolocationByIpAddress(context.Background(), "192.168.1.1")

	assert.NoError(t, err)
	assert.Equal(t, expectedGeolocation, geolocation)
}

func TestGetGeolocationByIpAddress_Failure(t *testing.T) {
	wrongIpAddress := "123.23.123.1"
	expectedErr := model.NewGeoLocationNotFoundException(wrongIpAddress)

	mockMongoDB.
		On("GetGeolocationByIpAddress", mock.Anything, wrongIpAddress).
		Return(nil, mongo.ErrNoDocuments)

	geolocation, err := service.GetGeolocationByIpAddress(context.Background(), wrongIpAddress)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), expectedErr.Error())
	assert.Nil(t, geolocation)
}

func TestGetGeolocations_Success(t *testing.T) {
	expectedPageResponse := &model.PageResponse{}
	mockMongoDB.On("GetPageableGeolocations", mock.Anything, 1, 10).Return(expectedPageResponse, nil).Times(1)

	pageResponse, err := service.GetGeolocations(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, expectedPageResponse, pageResponse)
}

func TestGetGeolocations_Failure(t *testing.T) {
	expectedError := errors.New("failed to retrieve geolocations")

	mockMongoDB.
		On("GetPageableGeolocations", mock.Anything, 1, 10).
		Return(nil, expectedError).Times(1)

	pageResponse, err := service.GetGeolocations(context.Background(), 1, 10)

	assert.Error(t, err)
	assert.Nil(t, pageResponse)
	assert.Equal(t, expectedError, err)
}
