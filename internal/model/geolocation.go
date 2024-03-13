package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Geolocation struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	IpAddress    string             `json:"ipAddress" bson:"ip_address"`
	CountryCode  string             `json:"countryCode" bson:"country_code"`
	Country      string             `json:"country" bson:"country"`
	City         string             `json:"city" bson:"city"`
	Latitude     float64            `json:"latitude" bson:"latitude"`
	Longitude    float64            `json:"longitude" bson:"longitude"`
	MysteryValue int64              `json:"mysteryValue" bson:"mystery_value"`
	CreatedAt    time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updated_at"`
}
