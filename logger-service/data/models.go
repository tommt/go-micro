package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client

func New(client *mongo.Client) Models {
	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID       string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string    `bson:"name" json:"name"`
	Data     string    `bson:"data" json:"data"`
	CreateAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		ID:       entry.ID,
		Name:     entry.Name,
		Data:     entry.Data,
		CreateAt: entry.CreateAt,
		UpdateAt: entry.UpdateAt,
	})
	if err != nil {
		log.Println("Error inserting log entry: ", err)
		return err
	}

	return nil
}

func (l *LogEntry) FindAll() ([]*LogEntry, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		log.Println("Error getting all log entries: ", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var entries []*LogEntry
	for cursor.Next(ctx) {
		var entry LogEntry
		err := cursor.Decode(&entry)
		if err != nil {
			log.Println("Error decoding log entry: ", err)
			return nil, err
		}
		entries = append(entries, &entry)
	}

	return entries, nil
}
