package storage

import (
	"LogForge/internal/model"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var LogCollection *mongo.Collection

func InitMongo() {
	client, err := mongo.NewClient(
		options.Client().ApplyURI("mongodb://localhost:27017"),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	LogCollection = client.Database("logforge").Collection("logs")

	log.Println("MongoDB connected")

	createIndexes()
}

func SaveLog(entry model.LogEntry) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := LogCollection.InsertOne(ctx, entry)
	if err != nil {
		log.Println("Mongo insert error:", err)
	}
}

func QueryLogs(service, level string, limit int64, from, to time.Time) ([]model.LogEntry, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}

	if service != "" {
		filter["service"] = service
	}

	if level != "" {
		filter["level"] = level
	}

	// Time range filter
	if !from.IsZero() || !to.IsZero() {
		timeFilter := bson.M{}

		if !from.IsZero() {
			timeFilter["$gte"] = from
		}
		if !to.IsZero() {
			timeFilter["$lte"] = to
		}

		filter["timestamp"] = timeFilter
	}

	opts := options.Find().
		SetLimit(limit).
		SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cursor, err := LogCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.LogEntry

	for cursor.Next(ctx) {
		var entry model.LogEntry
		if err := cursor.Decode(&entry); err != nil {
			return nil, err
		}
		results = append(results, entry)
	}

	return results, nil
}

func createIndexes() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		// Time-based queries
		{
			Keys: bson.D{{Key: "timestamp", Value: -1}},
		},

		// Service filter
		{
			Keys: bson.D{{Key: "service", Value: 1}},
		},

		// Level filter
		{
			Keys: bson.D{{Key: "level", Value: 1}},
		},

		// Compound index (most realistic queries)
		{
			Keys: bson.D{
				{Key: "service", Value: 1},
				{Key: "level", Value: 1},
				{Key: "timestamp", Value: -1},
			},
		},
	}

	_, err := LogCollection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Println("Index creation error:", err)
	} else {
		log.Println("Indexes created")
	}
}
