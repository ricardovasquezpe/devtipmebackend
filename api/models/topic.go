package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Topic struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Total     int                `json:"total" bson:"total"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func (t *Topic) Prepare() {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}

func (t *Topic) Validate() error {
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("Title is required")
	}
	return nil
}

func (t *Topic) SaveTopic(database *mongo.Database) (*Topic, error) {
	collection := database.Collection("topics")
	result, err := collection.InsertOne(context.TODO(), t)
	if err != nil {
		return &Topic{}, err
	}

	t.ID = result.InsertedID.(primitive.ObjectID)
	return t, nil
}

func FindByTitleAndIncrease(database *mongo.Database, title string) error {
	collection := database.Collection("topics")
	update := bson.M{
		"$inc": bson.M{"total": 1},
	}
	singleResult := collection.FindOneAndUpdate(context.TODO(), bson.M{"title": bson.M{"$regex": "^" + title + "$", "$options": "i"}}, update)
	if singleResult.Err() != nil {
		return singleResult.Err()
	}

	return nil
}

func GetTopicsLimited(database *mongo.Database, limit int64) ([]Topic, error) {
	var topics []Topic = []Topic{}
	collection := database.Collection("topics")
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{"total", -1}})

	top, err := collection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return []Topic{}, err
	}

	for top.Next(context.TODO()) {
		var topic Topic
		err = top.Decode(&topic)
		if err != nil {
			return []Topic{}, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
}
