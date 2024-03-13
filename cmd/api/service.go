package main

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"

	"geolocation-service/config"
	"geolocation-service/internal/model"
	"geolocation-service/pkg/logger"
)

type IGeolocationService interface {
	GetGeolocationByIpAddress(ctx context.Context, ipAddress string) (*model.Geolocation, error)
	GetGeolocations(ctx context.Context, offset int, limit int) (*model.PageResponse, error)
}

type GeolocationService struct {
	cfg config.Config
}

func NewGeolocationService(cfg *config.Config) IGeolocationService {
	return &GeolocationService{
		cfg: *cfg,
	}
}

func (s GeolocationService) GetGeolocationByIpAddress(ctx context.Context, ipAddress string) (*model.Geolocation, error) {
	log := logger.GetLogger(ctx)

	geolocation, err := s.cfg.MongoDB.GetGeolocationByIpAddress(ctx, ipAddress)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Errorf("ip address: %s not found", err.Error())
		return nil, model.NewGeoLocationNotFoundException(ipAddress)
	} else if err != nil {
		log.Errorf("Failed to get geolocation by ip address: %s", err.Error())
		return nil, err
	}

	return geolocation, nil
}

func (s GeolocationService) GetGeolocations(ctx context.Context, offset int, limit int) (*model.PageResponse, error) {
	log := logger.GetLogger(ctx)

	pageResponse, err := s.cfg.MongoDB.GetPageableGeolocations(ctx, offset, limit)
	if err != nil {
		log.Errorf("Failed to get pageable geolocations: %s", err.Error())
		return nil, err
	}

	return pageResponse, nil
}
