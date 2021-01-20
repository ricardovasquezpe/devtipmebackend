// user.go

package models

import (
	"context"
	"errors"
	"strings"

	_ "fmt"

	"github.com/badoux/checkmail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	Email        string `json:"email,omitempty" bson:"email,omitempty"`
	FirstName    string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName     string `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Password     string `json:"password,omitempty" bson:"password,omitempty"`
	ProfileImage string `json:"profileimage,omitempty" bson:"profileimage,omitempty"`
}

// HashPassword hashes password from user input
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14 is the cost for hashing the password.
	return string(bytes), err
}

// CheckPasswordHash checks password hash and password from user input if they match
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("password incorrect")
	}
	return nil
}

// BeforeSave hashes user password
func (u *User) BeforeSave() error {
	password := strings.TrimSpace(u.Password)
	hashedpassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.Password = string(hashedpassword)
	return nil
}

// Prepare strips user input of any white spaces
func (u *User) Prepare() {
	u.Email = strings.TrimSpace(u.Email)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.ProfileImage = strings.TrimSpace(u.ProfileImage)
}

// Validate user input
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
	default: // this is for creating a user, where all fields are required
		if u.FirstName == "" {
			return errors.New("FirstName is required")
		}
		if u.LastName == "" {
			return errors.New("LastName is required")
		}
		if u.Email == "" {
			return errors.New("Email is required")
		}
		if u.Password == "" {
			return errors.New("Password is required")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func GetUsers(database *mongo.Database) ([]User, error) {
	var users []User = []User{}
	collection := database.Collection("users1")
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
	_, err := collection.InsertOne(context.TODO(), u)
	if err != nil {
		return &User{}, err
	}

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
