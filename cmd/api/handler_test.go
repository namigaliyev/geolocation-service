package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"geolocation-service/config"
	"geolocation-service/internal/model"
	"geolocation-service/mocks/api"
)

var geolocationService = new(api.GeolocationServiceMock)

const pathPrefix = "/v1/geolocation-service"

func makeRequest(method, url string, apply func()) *httptest.ResponseRecorder {
	cfg := &config.Config{
		PathPrefix: pathPrefix,
	}
	mux := chi.NewMux()
	NewGeolocationHandler(*cfg, mux, geolocationService)
	req := httptest.NewRequest(method, cfg.PathPrefix+url, nil)
	apply()
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestGeolocationHandler_HandleGetGeolocationByIpAddress_Success(t *testing.T) {

	locationID := primitive.NewObjectID()
	ipAddress := "192.168.1.1"
	expectedResponse := &model.Geolocation{
		City:         "Baku",
		CountryCode:  "AZ",
		ID:           locationID,
		IpAddress:    ipAddress,
		Latitude:     36.6532,
		Longitude:    86.5035,
		MysteryValue: 2421,
	}

	applyFunction := func() {
		geolocationService.
			On("GetGeolocationByIpAddress",
				mock.Anything, ipAddress).
			Once().
			Return(expectedResponse, nil)
	}

	rec := makeRequest(http.MethodGet,
		"/geolocations/"+ipAddress,
		applyFunction)

	actualResponse := new(model.Geolocation)
	err := json.NewDecoder(rec.Body).Decode(actualResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedResponse, actualResponse)
	fmt.Println("success case")
}

func TestGeolocationHandler_HandleGetGeolocationByIpAddress_Failure_Invalid_Ip_Address(t *testing.T) {

	invalidIpAddress := "123312"
	expectedErr := model.NewGeoLocationNotFoundException(invalidIpAddress)

	applyFunction := func() {
		geolocationService.
			On("GetGeolocationByIpAddress",
				mock.Anything, invalidIpAddress).
			Once().
			Return(nil, expectedErr)
	}

	rec := makeRequest(http.MethodGet,
		"/geolocations/"+invalidIpAddress,
		applyFunction)

	errorBody := new(model.Exception)
	err := json.NewDecoder(rec.Body).Decode(errorBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, errorBody.HttpCode)
	fmt.Println("invalid ip address case")
}

func TestGeolocationHandler_HandleGetGeolocations(t *testing.T) {

	pageResponse := &model.PageResponse{
		Items: []interface{}{
			map[string]interface{}{
				"city":         "Baku",
				"countryCode":  "AZ",
				"createdAt":    "2023-11-04T20:34:58.651387237Z",
				"ipAddress":    "192.168.1.1",
				"latitude":     36.6532,
				"longitude":    86.5035,
				"mysteryValue": 2421.12321321,
			},
			map[string]interface{}{
				"city":         "Istanbul",
				"countryCode":  "TR",
				"createdAt":    "2023-11-04T20:34:58.651387237Z",
				"ipAddress":    "192.168.1.1",
				"latitude":     44.23124,
				"longitude":    55.5464,
				"mysteryValue": 789.23123,
			},
		},
		HasNextPage:    false,
		TotalPageCount: 1,
	}

	applyFunction := func() {
		geolocationService.
			On("GetGeolocations",
				mock.Anything, mock.Anything, mock.Anything).
			Once().
			Return(pageResponse, nil)
	}

	rec := makeRequest(http.MethodGet,
		"/geolocations/?offset=1&limit=10",
		applyFunction)

	actualResponse := new(model.PageResponse)
	err := json.NewDecoder(rec.Body).Decode(actualResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, pageResponse, actualResponse)

}
