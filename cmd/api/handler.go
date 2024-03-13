package main

import (
	"encoding/json"
	"fmt"
	"geolocation-service/internal/handler"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"geolocation-service/config"
	"geolocation-service/internal/model"
	"geolocation-service/pkg/utils"
)

type GeolocationHandler struct {
	geolocationService IGeolocationService
}

const (
	baseGroup = "/geolocations"

	getGeolocationByIpAddress = "/{ipAddress}"
)

func NewGeolocationHandler(cfg config.Config, router *chi.Mux, geolocationService IGeolocationService) *GeolocationHandler {
	h := &GeolocationHandler{
		geolocationService: geolocationService,
	}

	router.Route(cfg.PathPrefix+baseGroup, func(r chi.Router) {
		r.Get("/", handler.HandleError(h.HandleGetGeolocations))
		r.Get(getGeolocationByIpAddress, handler.HandleError(h.HandleGetGeolocationByIpAddress))
	})

	return h
}

// HandleGetGeolocations godoc
// @Tags geolocation-handler
// @Summary Endpoint for getting geolocations
// @Description Endpoint for getting pageable geolocations data
// @Param offset query int true "Page number (offset)"
// @Param limit query int true "Number of items per page (limit)"
// @Accept json
// @Produce json
// @Success 200 {object} model.PageResponse
// @Failure 404 {object} model.Exception
// @Failure 500 {object} model.Exception
// @Router  /api/geolocation-service/geolocations [get]
func (n *GeolocationHandler) HandleGetGeolocations(w http.ResponseWriter, r *http.Request) error {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	result, err := n.geolocationService.GetGeolocations(r.Context(), offset, limit)
	if err != nil {
		http.Error(w, "Failed to get geolocations", http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	return err
}

// HandleGetGeolocationByIpAddress godoc
// @Tags geolocation-handler
// @Summary Endpoint for getting geolocation
// @Description Endpoint for getting geolocation  by given ip address
// @Param ipAddress path string true "ip address"
// @Accept json
// @Produce json
// @Success 200 {object} model.Geolocation
// @Failure 404 {object} model.Exception
// @Failure 500 {object} model.Exception
// @Router  /api/geolocation-service/geolocations/{ipAddress} [get]
func (n *GeolocationHandler) HandleGetGeolocationByIpAddress(w http.ResponseWriter, r *http.Request) error {
	ipAddress := chi.URLParam(r, "ipAddress")

	if !utils.IsValidIP(ipAddress) {
		fmt.Printf("Invalid ip address: %s", ipAddress)
		return model.NewInvalidIPAddressException(ipAddress)
	}

	result, err := n.geolocationService.GetGeolocationByIpAddress(r.Context(), ipAddress)

	if err == nil {
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(result)
	}
	return err
}
