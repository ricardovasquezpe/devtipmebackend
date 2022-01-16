package models

import (
	"context"
	"devtipmebackend/utils"
	"errors"
	_ "fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const solutionStatusEnabled int = 1
const solutionStatusDisabled int = 0

type Solution struct {
	ID        *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string              `json:"title" bson:"title"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt" bson:"updatedAt"`
	Status    int                 `json:"status" bson:"status"`
	Content   []Content           `json:"content" bson:"content"`
	UserId    *primitive.ObjectID `json:"userId" bson:"userId"`
	Topics    []string            `json:"topics" bson:"topics"`
}

type SolutionResponse struct {
	ID          *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string              `json:"title" bson:"title"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	Status      *int                `json:"status" bson:"status"`
	Content     []Content           `json:"content" bson:"content"`
	EncriptedId *string             `json:"encriptedId" bson:"encriptedId"`
}

type Content struct {
	Type    int    `json:"type" bson:"type"`
	Content string `json:"content" bson:"content"`
	Order   int    `json:"order" bson:"order"`
}

func (s *Solution) Prepare() {
	s.Title = strings.TrimSpace(s.Title)
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	s.Status = solutionStatusEnabled
}

func (s *Solution) Validate() error {
	if strings.TrimSpace(s.Title) == "" {
		return errors.New("Title is required")
	}
	if s.CreatedAt.IsZero() {
		return errors.New("CreatedAt is required")
	}
	if s.UpdatedAt.IsZero() {
		return errors.New("UpdatedAt is required")
	}
	if len(s.Content) == 0 {
		return errors.New("Content is required")
	}
	if len(s.Topics) == 0 {
		return errors.New("Topics is required")
	}

	for index, element := range s.Content {
		if element.Type == 0 {
			return errors.New("Type is required in Content index " + strconv.Itoa(index))
		}
		if element.Content == "" {
			return errors.New("Content is required in Content index " + strconv.Itoa(index))
		}
		if element.Order == 0 {
			return errors.New("Order is required in Content index " + strconv.Itoa(index))
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

	idGenerated := result.InsertedID.(primitive.ObjectID)
	s.ID = &idGenerated
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

	solution.ID = nil
	//solution.UserId = nil

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

func FindAllSolutions(database *mongo.Database, text string, limit int64, offset int64, topic string) ([]SolutionResponse, error) {
	var solutions []SolutionResponse = []SolutionResponse{}
	collection := database.Collection("solutions")

	opts := options.Find()
	opts.SetSort(bson.D{{"createdAt", -1}})
	opts.SetLimit(limit)
	opts.SetSkip(offset)
	opts.SetProjection(bson.M{"title": 1, "content": 1, "_id": 1, "createdAt": 1})

	/*query := bson.M{
		"$text": bson.M{
			"$search": text,
		},
	}*/

	query := bson.M{"$and": []interface{}{
		bson.M{"$or": []interface{}{
			bson.M{"title": bson.M{"$regex": text, "$options": "im"}},
			bson.M{"content.content": bson.M{"$regex": text, "$options": "im"}},
		}},
		bson.M{"topics": bson.M{"$regex": topic, "$options": "im"}},
		bson.M{"status": solutionStatusEnabled},
	}}

	sol, err := collection.Find(context.TODO(), query, opts)
	if err != nil {
		return []SolutionResponse{}, err
	}

	for sol.Next(context.TODO()) {
		var solution SolutionResponse
		err = sol.Decode(&solution)
		if err != nil {
			return []SolutionResponse{}, err
		}

		stringEncriptedId, err := utils.Encrypt([]byte(solution.ID.Hex()), os.Getenv("SECRET"))
		if err != nil {
			return []SolutionResponse{}, err
		}

		solution.EncriptedId = &stringEncriptedId
		solution.ID = nil
		solution.Status = nil
		solutions = append(solutions, solution)
	}
	return solutions, nil
}

func GetSolutionsByUserId(database *mongo.Database, userId string) ([]SolutionResponse, error) {
	var solutions []SolutionResponse = []SolutionResponse{}
	collection := database.Collection("solutions")
	opts := options.Find()
	opts.SetSort(bson.D{{"createdAt", -1}})
	opts.SetProjection(bson.M{"title": 1, "content": 1, "_id": 1, "createdAt": 1, "status": 1})

	docID, _ := primitive.ObjectIDFromHex(userId)
	query := bson.M{"userId": docID}

	sol, err := collection.Find(context.TODO(), query, opts)
	if err != nil {
		return []SolutionResponse{}, err
	}

	for sol.Next(context.TODO()) {
		var solution SolutionResponse
		err = sol.Decode(&solution)
		if err != nil {
			return []SolutionResponse{}, err
		}

		stringEncriptedId, err := utils.Encrypt([]byte(solution.ID.Hex()), os.Getenv("SECRET"))
		if err != nil {
			return []SolutionResponse{}, err
		}

		solution.EncriptedId = &stringEncriptedId
		solution.ID = nil
		solutions = append(solutions, solution)
	}
	return solutions, nil
}

func UpdateSolutionStatus(database *mongo.Database, solutionId string, status float64) error {
	collection := database.Collection("solutions")
	primitiveId, _ := primitive.ObjectIDFromHex(solutionId)
	filter := bson.M{"_id": primitiveId}
	update := bson.M{
		"$set": bson.M{"status": status, "updatedAt": time.Now()},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
