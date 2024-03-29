package models

import (
	"context"
	"devtipmebackend/utils"
	"errors"
	"os"
	"strings"
	"time"

	_ "fmt"

	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const userStatusEnabled int = 1
const userStatusDisabled int = 0

type User struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email string             `json:"email,omitempty" bson:"email,omitempty"`
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	//LastName     string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Password     string    `json:"password,omitempty" bson:"password,omitempty"`
	ProfileImage string    `json:"profileimage,omitempty" bson:"profileimage,omitempty"`
	Status       int       `json:"status" bson:"status"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt"`
}

type UserResponse struct {
	ID          *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email       string              `json:"email,omitempty" bson:"email,omitempty"`
	Name        string              `json:"name,omitempty" bson:"name,omitempty"`
	EncriptedId *string             `json:"encriptedId" bson:"encriptedId"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("password incorrect")
	}
	return nil
}

func (u *User) BeforeSave() error {
	password := strings.TrimSpace(u.Password)
	hashedpassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.Password = string(hashedpassword)
	u.Status = 0
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) Prepare() {
	u.Email = strings.TrimSpace(u.Email)
	u.Name = strings.TrimSpace(u.Name)
	//u.LastName = strings.TrimSpace(u.LastName)
	u.ProfileImage = strings.TrimSpace(u.ProfileImage)
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Email == "" {
			return errors.New("Email is required")
		}
		if u.Password == "" {
			return errors.New("Password is required")
		}
		return nil
	default:
		if u.Name == "" {
			return errors.New("Name is required")
		}
		if u.Email == "" {
			return errors.New("Email is required")
		}
		if u.Password == "" {
			return errors.New("Password is required")
		}
		if len(u.Password) < 6 {
			return errors.New("Password should have more than 6 characters")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func GetUsers(database *mongo.Database) ([]User, error) {
	var users []User = []User{}
	collection := database.Collection("users")
	usr, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return []User{}, err
	}

	for usr.Next(context.TODO()) {
		var user User
		err = usr.Decode(&user)
		if err != nil {
			return []User{}, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *User) SaveUser(database *mongo.Database) (*User, error) {
	collection := database.Collection("users")

	result, err := collection.InsertOne(context.TODO(), u)
	if err != nil {
		return &User{}, err
	}
	u.ID = result.InsertedID.(primitive.ObjectID)

	return u, nil
}

func DeleteUser(id string, database *mongo.Database) error {
	collection := database.Collection("users")
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": idPrimitive})
	_ = res
	if err != nil {
		return err
	}
	return nil
}

func (u *User) UpdateUser(id string, database *mongo.Database) (*User, error) {
	collection := database.Collection("users")
	primitiveId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": primitiveId}
	update := bson.M{
		"$set": u,
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) GetUserByEmail(database *mongo.Database) (*User, error) {
	user := &User{}
	collection := database.Collection("users")

	filter := bson.M{"email": u.Email}
	err := collection.FindOne(context.TODO(), filter).Decode(user)

	if err != nil {
		return &User{}, err
	}

	return user, nil
}

func GetUserById(database *mongo.Database, id string) (*User, error) {
	user := &User{}
	collection := database.Collection("users")

	docID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": docID}
	err := collection.FindOne(context.TODO(), filter).Decode(user)

	if err != nil {
		return &User{}, err
	}

	return user, nil
}

func GetUserByIdExternal(database *mongo.Database, id string) (*UserResponse, error) {
	user := &UserResponse{}
	collection := database.Collection("users")

	docID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": docID}
	opts := options.FindOne()
	opts.SetProjection(bson.M{"email": 1, "name": 1, "_id": 1})

	err := collection.FindOne(context.TODO(), query, opts).Decode(user)

	if err != nil {
		return &UserResponse{}, err
	}

	stringEncriptedId, err := utils.Encrypt([]byte(user.ID.Hex()), os.Getenv("SECRET"))
	if err != nil {
		return &UserResponse{}, err
	}

	user.EncriptedId = &stringEncriptedId
	user.ID = nil

	return user, nil
}
