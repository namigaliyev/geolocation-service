package main

import (
	"context"
	"geolocation-service/pkg/logger"
	"sync"

	"geolocation-service/config"
)

const addr = ":8002"

func main() {
	ctx := context.Background()
	log := logger.GetLogger(ctx)
	var wg sync.WaitGroup

	log.Info("Starting import-service")

	cfg, err := config.NewConfig(addr)
	if err != nil {
		log.Fatalln("could not setup config, ", err)
	}
	defer cfg.MongoDB.Close()

	importerService := NewImporterService(cfg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := importerService.Start()
		if err != nil {
			log.Errorf("Unable to start importer service: %s", err)
		}
	}()

	wg.Wait()

	log.Info("Terminating import-service")
}
