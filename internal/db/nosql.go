package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"geolocation-service/internal/common"
	"geolocation-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"math"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	geolocationsCollection = "geolocations"
	metricsCollection      = "metrics"
)

type NoSQLDb interface {
	// CheckHealth returns the status of the store.
	CheckHealth(ctx context.Context) bool

	SaveGeolocations(ctx context.Context, locations []*model.Geolocation) error

	GetGeolocations(ctx context.Context) ([]*model.Geolocation, error)

	GetGeolocationByIpAddress(ctx context.Context, ipAddress string) (*model.Geolocation, error)

	GetPageableGeolocations(ctx context.Context, offset int, limit int) (*model.PageResponse, error)

	DeleteAllGeolocations(ctx context.Context) error

	SaveMetric(ctx context.Context, metrics *model.Metric) error

	GetMetricByProcessID(ctx context.Context, processID string) (*model.Metric, error)

	UpdateMetric(ctx context.Context, metric *model.Metric) error

	DeleteMetricByProcessID(ctx context.Context, processID string) error

	// Close terminates any MongoDB connections gracefully.
	Close() error
}

type MongoDB struct {
	Client *mongo.Client

	DB *mongo.Database
}

// NewMongoDB returns new a MongoDB client.
func NewMongoDB(uri, dbName, certPath string) (*MongoDB, error) {

	clientOptions := options.Client().ApplyURI(uri)

	if certPath != "" {
		c, err := getCustomTLSConfig(certPath)
		if err != nil {
			return nil, err
		}

		clientOptions.SetTLSConfig(c)
	}

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	mongoDB := &MongoDB{
		Client: client,
		DB:     client.Database(dbName),
	}

	return mongoDB, nil
}

func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	certs, err := os.ReadFile(filepath.Clean(caFile))

	if err != nil {
		return nil, err
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

	if !ok {
		return nil, errors.New("failed parsing pem file")
	}

	return tlsConfig, nil
}

// Close terminates any MongoDB connections gracefully.
func (mongoDB *MongoDB) Close() error {
	return mongoDB.Client.Disconnect(context.TODO())
}

// CheckHealth returns the status of the store.
func (mongoDB *MongoDB) CheckHealth(ctx context.Context) bool {
	err := mongoDB.Client.Ping(ctx, readpref.Primary())

	return err == nil
}

func (mongoDB *MongoDB) SaveGeolocations(ctx context.Context, users []*model.Geolocation) error {
	coll := mongoDB.DB.Collection(geolocationsCollection)

	var insertModels []mongo.WriteModel
	for _, user := range users {
		insertModel := mongo.NewInsertOneModel().SetDocument(user)
		insertModels = append(insertModels, insertModel)
	}

	opts := options.BulkWrite().SetOrdered(false)
	result, err := coll.BulkWrite(ctx, insertModels, opts)
	if err != nil {
		return err
	}

	if int(result.InsertedCount) != len(users) {
		return errors.New("some locations weren't inserted successfully")
	}

	return nil
}

func (mongoDB *MongoDB) GetGeolocations(ctx context.Context) ([]*model.Geolocation, error) {
	matchStage := bson.D{{"$match", bson.D{{}}}}
	eventTypes := make([]*model.Geolocation, 0)

	coll := mongoDB.DB.Collection(geolocationsCollection)
	cur, err := coll.Aggregate(ctx, mongo.Pipeline{matchStage})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &eventTypes)
	if err != nil {
		return nil, err
	}

	return eventTypes, nil
}

func (mongoDB *MongoDB) GetGeolocationByIpAddress(ctx context.Context, ipAddress string) (*model.Geolocation, error) {
	coll := mongoDB.DB.Collection(geolocationsCollection)

	var geolocation model.Geolocation
	err := coll.FindOne(ctx, bson.M{"ip_address": ipAddress}).Decode(&geolocation)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, mongo.ErrNoDocuments
	}

	return &geolocation, nil
}

func (mongoDB *MongoDB) GetPageableGeolocations(ctx context.Context, offset int, limit int) (*model.PageResponse, error) {
	matchStage := bson.D{{"$match", bson.D{{}}}}
	coll := mongoDB.DB.Collection(geolocationsCollection)

	totalCount, err := coll.CountDocuments(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	skip := (offset - 1) * limit
	skipStage := bson.D{{"$skip", skip}}
	limitStage := bson.D{{"$limit", limit}}

	users := make([]*model.Geolocation, 0)
	cur, err := coll.Aggregate(ctx, mongo.Pipeline{matchStage, skipStage, limitStage})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &users)
	if err != nil {
		return nil, err
	}

	totalPageCount := int(math.Ceil(float64(totalCount) / float64(limit)))
	hasNextPage := offset*limit < int(totalCount)

	return &model.PageResponse{
		Items:          users,
		HasNextPage:    hasNextPage,
		TotalPageCount: totalPageCount,
	}, nil
}

func (mongoDB *MongoDB) DeleteAllGeolocations(ctx context.Context) error {
	coll := mongoDB.DB.Collection(geolocationsCollection)
	_, err := coll.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}
	return nil
}

func (mongoDB *MongoDB) SaveMetric(ctx context.Context, metrics *model.Metric) error {
	collection := mongoDB.DB.Collection(metricsCollection)
	_, err := collection.InsertOne(ctx, metrics)

	return err
}

func (mongoDB *MongoDB) GetMetricByProcessID(ctx context.Context, processID string) (*model.Metric, error) {
	coll := mongoDB.DB.Collection(metricsCollection)

	var metric model.Metric
	err := coll.FindOne(ctx, bson.M{"process_id": processID}).Decode(&metric)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, mongo.ErrNoDocuments
	}

	return &metric, nil
}

func (mongoDB *MongoDB) UpdateMetric(ctx context.Context, metric *model.Metric) error {
	coll := mongoDB.DB.Collection(metricsCollection)

	filter := bson.M{"_id": metric.ID}
	update := bson.M{
		"$set": bson.M{
			"updated_at": metric.UpdatedAt,
		},
		"$inc": bson.M{
			"checkpoint_off_set":   metric.CheckpointOffSet,
			"accepted_records":     metric.AcceptedRecords,
			"discarded_records":    metric.DiscardedRecords,
			"duplicate_records":    metric.DuplicateRecords,
			"malformed_ips":        metric.MalformedIPs,
			"invalid_latitude":     metric.InvalidLatitude,
			"invalid_longitude":    metric.InvalidLongitude,
			"elapsed_time_seconds": metric.ElapsedTimeSeconds,
		},
	}

	result := coll.FindOneAndUpdate(ctx, filter, update)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return common.ErrNoRecordToUpdate
		}
		return result.Err()
	}

	return nil
}

func (mongoDB *MongoDB) DeleteMetricByProcessID(ctx context.Context, processID string) error {
	coll := mongoDB.DB.Collection(metricsCollection)

	// Delete the document with the given process ID
	_, err := coll.DeleteOne(ctx, bson.M{"process_id": processID})
	if err != nil {
		return err
	}

	return nil
}
