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

type Comment struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Comment    string             `json:"comment" bson:"comment"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt"`
	SolutionId primitive.ObjectID `json:"solutionId" bson:"solutionId"`
	UserId     primitive.ObjectID `json:"userId" bson:"userId"`
}

type CommentResponse struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Comment   string             `json:"comment" bson:"comment"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	UserData  User               `json:"userData,omitempty" bson:"userData,omitempty"`
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

func FindAllComments(database *mongo.Database, solutionId string) ([]CommentResponse, error) {
	var comments []CommentResponse = []CommentResponse{}
	collection := database.Collection("comments")
	docID, _ := primitive.ObjectIDFromHex(solutionId)
	/*query := bson.M{"solutionId": docID}
	com, err := collection.Find(context.TODO(), query)*/

	matchStage := bson.D{{"$match", bson.D{{"solutionId", docID}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "users"}, {"localField", "userId"}, {"foreignField", "_id"}, {"as", "userData"}}}}
	unwind := bson.D{{"$unwind", "$userData"}}
	project := bson.D{{"$project", bson.D{{"userData.password", 0}, {"userData._id", 0}, {"userData.createdAt", 0}, {"userData.updatedAt", 0}, {"userData.status", 0}}}}
	sort := bson.D{{"$sort", bson.D{{"createdAt", -1}}}}

	opts := options.Aggregate()
	com, err := collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, lookupStage, unwind, project, sort}, opts)

	if err != nil {
		return []CommentResponse{}, nil
	}

	for com.Next(context.TODO()) {
		var comment CommentResponse
		err = com.Decode(&comment)
		if err != nil {
			return []CommentResponse{}, nil
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
