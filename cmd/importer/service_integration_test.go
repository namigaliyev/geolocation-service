package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"geolocation-service/config"
)

func TestImportService(t *testing.T) {
	var wg sync.WaitGroup

	cfg, err := config.NewConfig(addr)
	if err != nil {
		log.Fatalln("could not setup config, ", err)
	}
	defer cfg.MongoDB.Close()
	cfg.DataCheckpointLocation = "../../data_dump.csv"

	importerService := NewImporterService(cfg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := importerService.Start()
		if err != nil {
			fmt.Println(err)
		}
	}()

	wg.Wait()

	metric, err := cfg.MongoDB.GetMetricByProcessID(context.Background(), cfg.DataCheckpointLocation)
	if err != nil {
		t.Errorf("failed to get metric: %s", err.Error())
	}
	if metric == nil {
		t.Errorf("metric is nil")
	}

	t.Log("CheckpointOffset:", metric.CheckpointOffSet)
	t.Log("AcceptedRecords:", metric.AcceptedRecords)
	t.Log("DiscardedRecords:", metric.DiscardedRecords)
	t.Log("DuplicateRecords:", metric.DuplicateRecords)
	t.Log("InvalidLatitude:", metric.InvalidLatitude)
	t.Log("InvalidLongitude:", metric.InvalidLongitude)
	t.Log("MalformedIPs:", metric.MalformedIPs)
	t.Log("ElapsedTimeSeconds:", metric.ElapsedTimeSeconds)
	assert.Equal(t, metric.CheckpointOffSet, 1000000)
	assert.Equal(t, metric.AcceptedRecords, 866449)
	assert.Equal(t, metric.DiscardedRecords, 50451)
	assert.Equal(t, metric.DuplicateRecords, 50156)
	assert.Equal(t, metric.InvalidLatitude, 16402)
	assert.Equal(t, metric.InvalidLongitude, 0)
	assert.Equal(t, metric.MalformedIPs, 16542)

	// Clean up
	err = cfg.MongoDB.DeleteAllGeolocations(context.Background())
	if err != nil {
		t.Errorf("failed to delete geolocations: %s", err.Error())
	}
	err = cfg.MongoDB.DeleteMetricByProcessID(context.Background(), cfg.DataCheckpointLocation)
	if err != nil {
		t.Errorf("failed to delete metric: %s", err.Error())
	}

	t.Log("Success")
}
