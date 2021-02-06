package models

import (
	"context"
	"errors"
	_ "fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Solution struct {
	Title      string    `json:"title" bson:"title"`
	RegisterAt time.Time `json:"registerAt" bson:"registerAt"`
	UpdatedAt  time.Time `json:"updatedAt" bson:"updatedAt"`
	Tips       int       `json:"tips" bson:"tips"`
	Content    []Content `json:"content" bson:"content"`
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
	_, err := collection.InsertOne(context.TODO(), s)
	if err != nil {
		return &Solution{}, err
	}

	return s, nil
}
