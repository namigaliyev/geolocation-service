package config

import (
	"fmt"
	"os"

	"geolocation-service/internal/db"
)

type Config struct {
	Addr                   string
	PathPrefix             string
	MongoDB                db.NoSQLDb
	DataCheckpointLocation string
}

func NewConfig(addr string) (*Config, error) {
	cfg := new(Config)
	cfg.Addr = addr
	cfg.PathPrefix = os.Getenv("PATH_PREFIX")
	cfg.DataCheckpointLocation = os.Getenv("DATA_CHECKPOINT_LOCATION")

	m, err := getMongoDB()
	if err != nil {
		return nil, err
	}

	cfg.MongoDB = m

	return cfg, nil
}

func getMongoDB() (db.NoSQLDb, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	mongoDBName := os.Getenv("MONGODB_DATABASE_NAME")
	mongoCertPath := os.Getenv("MONGODB_CERT_PATH")

	mongoDB, err := db.NewMongoDB(mongoURI, mongoDBName, mongoCertPath)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %s", err)
	}

	return mongoDB, nil
}
