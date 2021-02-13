package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Comment struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Comment    string             `json:"comment" bson:"comment"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt"`
	SolutionId primitive.ObjectID `json:"solutionId" bson:"solutionId"`
	UserId     primitive.ObjectID `json:"userId" bson:"userId"`
}

func (c *Comment) Prepare() {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	c.Comment = strings.TrimSpace(c.Comment)
}

func (c *Comment) Validate() error {
	if c.Comment == "" {
		return errors.New("Comment is required")
	}
	if c.CreatedAt.IsZero() {
		return errors.New("CreatedAt is required")
	}
	if c.UpdatedAt.IsZero() {
		return errors.New("UpdatedAt is required")
	}
	return nil
}

func (c *Comment) SaveComment(database *mongo.Database) (*Comment, error) {
	collection := database.Collection("comments")
	result, err := collection.InsertOne(context.TODO(), c)
	if err != nil {
		return &Comment{}, err
	}

	c.ID = result.InsertedID.(primitive.ObjectID)
	return c, nil
}
