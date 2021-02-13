package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Tip struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Amount     float64            `json:"amount" bson:"amount"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt"`
	SolutionId primitive.ObjectID `json:"solutionId" bson:"solutionId"`
	UserId     primitive.ObjectID `json:"userId" bson:"userId"`
}

func (t *Tip) Prepare() {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}

func (t *Tip) Validate() error {
	/*if t.Amount > 0 {
		return errors.New("Amount is required")
	}*/
	if t.CreatedAt.IsZero() {
		return errors.New("CreatedAt is required")
	}
	if t.UpdatedAt.IsZero() {
		return errors.New("UpdatedAt is required")
	}
	/*if t.SolutionId != 0 {
		return errors.New("SolutionId is required")
	}*/
	return nil
}

func (t *Tip) SaveTip(database *mongo.Database) (*Tip, error) {
	collection := database.Collection("tips")
	result, err := collection.InsertOne(context.TODO(), t)
	if err != nil {
		return &Tip{}, err
	}

	t.ID = result.InsertedID.(primitive.ObjectID)
	return t, nil
}
