package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Metric struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id"`
	ProcessID          string             `json:"processID" bson:"process_id"`
	CheckpointOffSet   int                `json:"checkpointOffSet" bson:"checkpoint_off_set"`
	AcceptedRecords    int                `json:"acceptedRecords" bson:"accepted_records"`
	DiscardedRecords   int                `json:"discardedRecords" bson:"discarded_records"`
	DuplicateRecords   int                `json:"duplicateRecords" bson:"duplicate_records"`
	MalformedIPs       int                `json:"malformedIPs" bson:"malformed_ips"`
	InvalidLatitude    int                `json:"invalidLatitude" bson:"invalid_latitude"`
	InvalidLongitude   int                `json:"invalidLongitude" bson:"invalid_longitude"`
	ElapsedTimeSeconds float64            `json:"elapsedTimeSeconds" bson:"elapsed_time_seconds"`
	CreatedAt          time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updatedAt" bson:"updated_at"`
}
