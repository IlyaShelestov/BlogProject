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
	Username       string
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

func (m *UserModel) Authenticate(username, password string) (int, error) {
	filter := bson.D{{"username", username}}

	var user bson.M

	err := m.Collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	hashedPassword, ok := user["hashed_password"].(string)
	if !ok {
		return 0, ErrInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	id, ok := user["id"].(int)
	if !ok {
		return 0, ErrInvalidCredentials
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	filter := bson.D{{"id", id}}

	err := m.Collection.FindOne(context.Background(), filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
