package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserModelInterface interface {
	Insert(username, password string) error
	Authenticate(username, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (User, error)
	PasswordUpdate(id int, currentPassword, newPassword string) error
}

type User struct {
	ID             int
	username       string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	Collection *mongo.Collection
}

func NewUserModel(database *mongo.Database) *UserModel {
	return &UserModel{
		Collection: database.Collection("users"),
	}
}

func (m *UserModel) Insert(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user := bson.M{
		"username":        username,
		"hashed_password": hashedPassword,
		"created":         time.Now(),
	}

	_, err = m.Collection.InsertOne(context.Background(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrDuplicateUsername
		}
		return err
	}
	return nil
}
