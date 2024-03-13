package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"geolocation-service/config"
	"geolocation-service/internal/model"
	"geolocation-service/pkg/logger"
	"geolocation-service/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IImporterService interface {
	Start() error
}

type ImporterService struct {
	cfg config.Config
}

func NewImporterService(cfg *config.Config) IImporterService {
	return &ImporterService{
		cfg: *cfg,
	}
}

func (s ImporterService) Start() error {
	ctx := context.Background()
	log := logger.GetLogger(ctx)

	metric, err := s.getOrCreateMetric(ctx, s.cfg.DataCheckpointLocation)
	if err != nil {
		log.Fatalf("Failed to load or create metrics & checkpoint: %s", err)
	}

	csvFile, err := os.Open(s.cfg.DataCheckpointLocation)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %s", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	_, err = reader.Read()
	if err != nil {
		log.Fatalf("Failed to read header from CSV: %s", err)
	}

	batchSize := 5000
	seenIPs, _ := s.getAllSeenIPs(ctx)
	records := make(chan [][]string)
	results := make(chan error)
	var wg sync.WaitGroup
	var mu sync.RWMutex
	log.Infof("seen ip size %d", len(seenIPs))

	// Concurrent processing of records
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range records {
				// Process batch of records
				err := s.processBatchRecords(ctx, batch, metric.ID, &mu, seenIPs)
				if err != nil {
					results <- err
					return
				}
			}
		}()
	}
	// Read and process records in batches
	var batch [][]string
	currentOffset := 0
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Failed to read record from CSV: %s", err)
		}
		currentOffset++
		if metric.CheckpointOffSet >= currentOffset {
			continue
		}
		batch = append(batch, record)
		if len(batch) >= batchSize {
			select {
			case records <- batch:
				batch = nil
			case err := <-results:
				log.Fatalf("Error processing record: %s", err)
			}
		}
	}

	// Process remaining records
	if len(batch) > 0 {
		select {
		case records <- batch:
			// Wait for processing to finish
		case err = <-results:
			log.Fatalf("Error processing record: %s", err)
		}
	}

	close(records)

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Check for errors from goroutines
	for err := range results {
		if err != nil {
			log.Fatalf("Error processing record: %s", err)
		}
	}

	return nil
}

func (s ImporterService) processBatchRecords(
	ctx context.Context,
	batch [][]string,
	metricID primitive.ObjectID,
	mu *sync.RWMutex,
	seenIPs map[string]bool) error {

	// Initialize metrics for this batch
	batchMetrics := &model.Metric{}
	now := time.Now()
	log := logger.GetLogger(ctx)

	// Initialize a slice to hold geolocations for this batch
	geolocationBatch := make([]*model.Geolocation, 0, len(batch))

	for _, record := range batch {
		// Process the record
		geolocation, metrics, _ := s.processSingleRecord(record)

		if geolocation != nil {
			// Check for duplicate IP addresses
			mu.Lock()
			if _, ok := seenIPs[geolocation.IpAddress]; ok {
				batchMetrics.DuplicateRecords++
				mu.Unlock()
				continue
			}
			seenIPs[geolocation.IpAddress] = true
			mu.Unlock()

			geolocationBatch = append(geolocationBatch, geolocation)
		}

		// Update batch metrics
		batchMetrics.AcceptedRecords += metrics.AcceptedRecords
		batchMetrics.DiscardedRecords += metrics.DiscardedRecords
		batchMetrics.MalformedIPs += metrics.MalformedIPs
		batchMetrics.DuplicateRecords += metrics.DuplicateRecords
		batchMetrics.InvalidLatitude += metrics.InvalidLatitude
		batchMetrics.InvalidLongitude += metrics.InvalidLongitude
	}

	batchMetrics.ID = metricID
	batchMetrics.ElapsedTimeSeconds = time.Since(now).Seconds()
	batchMetrics.CheckpointOffSet = len(batch)
	log.Infof("Batch size: %d", len(geolocationBatch))
	if err := s.cfg.MongoDB.UpdateMetric(ctx, batchMetrics); err != nil {
		log.Errorf("failed to update metric: %v", err)
		return fmt.Errorf("failed to update metric: %w", err)
	}
	if err := s.cfg.MongoDB.SaveGeolocations(ctx, geolocationBatch); err != nil {
		log.Errorf("failed to save geolocations: %v", err)
		return fmt.Errorf("failed to save geolocations: %w", err)
	}

	return nil
}

func (s ImporterService) processSingleRecord(record []string) (*model.Geolocation, model.Metric, error) {
	metrics := model.Metric{}

	if len(record) < 7 || record[0] == "" || record[1] == "" || record[4] == "" || record[5] == "" {
		metrics.DiscardedRecords++
		return nil, metrics, fmt.Errorf("empty fields")
	}

	ip := record[0]
	countryCode := record[1]
	country := record[2]
	city := record[3]
	latitude, latErr := strconv.ParseFloat(record[4], 64)
	longitude, lonErr := strconv.ParseFloat(record[5], 64)
	mysteryValue, _ := strconv.ParseInt(record[6], 10, 64)

	// Validate IP address
	if !utils.IsValidIP(ip) {
		metrics.MalformedIPs++
		return nil, metrics, fmt.Errorf("malformed IP address")
	}

	// Validate latitude and longitude
	if latErr != nil || latitude < -90 || latitude > 90 {
		metrics.InvalidLatitude++
		return nil, metrics, fmt.Errorf("invalid latitude")
	}
	if lonErr != nil || longitude < -180 || longitude > 180 {
		metrics.InvalidLongitude++
		return nil, metrics, fmt.Errorf("invalid longitude")
	}

	geolocation := &model.Geolocation{
		ID:           primitive.NewObjectID(),
		IpAddress:    ip,
		CountryCode:  countryCode,
		Country:      country,
		City:         city,
		Latitude:     latitude,
		Longitude:    longitude,
		MysteryValue: mysteryValue,
	}

	metrics.AcceptedRecords++
	return geolocation, metrics, nil
}

func (s ImporterService) getOrCreateMetric(ctx context.Context, processID string) (*model.Metric, error) {
	now := time.Now()
	metric, _ := s.cfg.MongoDB.GetMetricByProcessID(ctx, processID)
	if metric != nil {
		return metric, nil
	}
	newMetric := &model.Metric{
		ID:        primitive.NewObjectID(),
		ProcessID: processID,
		CreatedAt: now,
	}
	if err := s.cfg.MongoDB.SaveMetric(ctx, newMetric); err != nil {
		return nil, fmt.Errorf("failed to save new metric: %w", err)
	}
	return newMetric, nil
}

func (s ImporterService) getAllSeenIPs(ctx context.Context) (map[string]bool, error) {
	seenIPS := make(map[string]bool)
	geolocations, err := s.cfg.MongoDB.GetGeolocations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations %s", err.Error())
	}
	for _, geolocation := range geolocations {
		seenIPS[geolocation.IpAddress] = true
	}
	return seenIPS, nil
}
