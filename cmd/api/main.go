package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"

	"geolocation-service/config"
	"geolocation-service/internal/handler"
	"geolocation-service/internal/middleware"
	"geolocation-service/pkg/logger"
)

const addr = ":8001"

func main() {
	ctx := context.Background()
	log := logger.GetLogger(ctx)
	cfg, err := config.NewConfig(addr)
	if err != nil {
		log.Fatalln("could not setup config, ", err)
	}
	defer cfg.MongoDB.Close()

	router := chi.NewRouter()
	router.Use(middleware.RequestParamsMiddleware)

	geolocationService := NewGeolocationService(cfg)

	handler.NewHealthHandler(router)
	NewGeolocationHandler(*cfg, router, geolocationService)

	log.Infof("Starting server at port %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
