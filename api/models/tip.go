package models

import (
	"context"
	"errors"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	if t.Amount == 0 {
		return errors.New("Amount is required")
	}
	if t.CreatedAt.IsZero() {
		return errors.New("CreatedAt is required")
	}
	if t.UpdatedAt.IsZero() {
		return errors.New("UpdatedAt is required")
	}
	/*if t.SolutionId.Hex() != "" {
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

func GetTotalTipBySolutionId(database *mongo.Database, solutionId string) (float64, error) {
	collection := database.Collection("tips")
	solID, _ := primitive.ObjectIDFromHex(solutionId)

	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":        "$solutionId",
				"sum_amount": bson.M{"$sum": "$amount"},
			},
		},
		{
			"$match": bson.M{"_id": solID},
		},
	}
	opts := options.Aggregate()
	cur, err := collection.Aggregate(context.TODO(), pipeline, opts)
	if err != nil {
		return 0.0, err
	}

	defer cur.Close(context.TODO())

	var doc []bson.M
	if err = cur.All(context.TODO(), &doc); err != nil {
		return 0, err
	}

	if len(doc) == 0 {
		return 0.0, nil
	}

	count := (doc[0]["sum_amount"]).(float64)
	count = math.Round(count*100) / 100

	return count, nil
}
