package models

import (
	"context"
	"errors"
	_ "fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Solution struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title      string             `json:"title" bson:"title"`
	RegisterAt time.Time          `json:"registerAt" bson:"registerAt"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt"`
	Content    []Content          `json:"content" bson:"content"`
}

type Content struct {
	Type    int    `json:"type" bson:"type"`
	Content string `json:"content" bson:"content"`
}

func (s *Solution) Prepare() {
	s.Title = strings.TrimSpace(s.Title)
	s.RegisterAt = time.Now()
	s.UpdatedAt = time.Now()
}

func (s *Solution) Validate() error {
	if s.Title == "" {
		return errors.New("Title is required")
	}
	if s.RegisterAt.IsZero() {
		return errors.New("RegisterAt is required")
	}
	if s.UpdatedAt.IsZero() {
		return errors.New("UpdatedAt is required")
	}
	if len(s.Content) == 0 {
		return errors.New("Content is required")
	}

	for index, element := range s.Content {
		if element.Type == 0 {
			return errors.New("Type is required in Content index " + strconv.Itoa(index))
		}

		if element.Content == "" {
			return errors.New("Content is required in Content index " + strconv.Itoa(index))
		}
	}

	return nil
}

func (s *Solution) SaveSolution(database *mongo.Database) (*Solution, error) {
	collection := database.Collection("solutions")
	result, err := collection.InsertOne(context.TODO(), s)
	if err != nil {
		return &Solution{}, err
	}

	s.ID = result.InsertedID.(primitive.ObjectID)
	return s, nil
}

func GetSolutionById(database *mongo.Database, id string) (*Solution, error) {
	solution := &Solution{}
	collection := database.Collection("solutions")

	docID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": docID}
	err := collection.FindOne(context.TODO(), filter).Decode(solution)

	if err != nil {
		return &Solution{}, err
	}

	return solution, nil
}

func GetAllSolutions(database *mongo.Database) ([]Solution, error) {
	var solutions []Solution = []Solution{}
	collection := database.Collection("solutions")
	sol, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return []Solution{}, err
	}

	for sol.Next(context.TODO()) {
		var solution Solution
		err = sol.Decode(&solution)
		if err != nil {
			return []Solution{}, err
		}
		solutions = append(solutions, solution)
	}
	return solutions, nil
}
