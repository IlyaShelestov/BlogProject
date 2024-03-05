package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(username, password string) error
	Authenticate(username, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (User, error)
	PasswordUpdate(id int, currentPassword, newPassword string) error
	ExistsByUsername(username string) (bool, error)
	IsAdmin(id int) (bool, error)
}

type User struct {
	ID             int
	Username       string
	HashedPassword string
	Created        time.Time
	IsAdmin        bool
}

type UserModel struct {
	Collection *mongo.Collection
}

func (m *UserModel) Insert(username, password string) error {
	lastID, err := m.getLastID()
	if err != nil {
		return err
	}

	newId := lastID + 1

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	hashedPassword := string(hashedPasswordBytes)

	user := bson.M{
		"id":              newId,
		"username":        username,
		"hashed_password": hashedPassword,
		"created":         time.Now(),
		"isAdmin":         false,
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
	filter := bson.D{{Key: "username", Value: username}}

	var user bson.M

	err := m.Collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
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
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	id := int(user["id"].(int32))

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	filter := bson.D{{Key: "id", Value: id}}

	err := m.Collection.FindOne(context.Background(), filter).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (m *UserModel) ExistsByUsername(username string) (bool, error) {
	filter := bson.D{{Key: "username", Value: username}}

	err := m.Collection.FindOne(context.Background(), filter).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (m *UserModel) Get(id int) (User, error) {
	filter := bson.D{{Key: "id", Value: id}}

	var user User

	err := m.Collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return User{}, ErrNoUserFound
		}
		return User{}, err
	}

	return user, nil
}

func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	var user bson.M
	filter := bson.D{{Key: "id", Value: id}}
	err := m.Collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNoUserFound
		}
		return err
	}

	hashedPassword, ok := user["hashed_password"].(string)
	if !ok {
		return ErrInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}
		return err
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "hashed_password", Value: newHashedPassword}}}}
	_, err = m.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) getLastID() (int32, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
	filter := bson.D{}
	var user bson.M
	err := m.Collection.FindOne(context.Background(), filter, opts).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, nil
		}
		return 0, err
	}
	lastID, ok := user["id"].(int32)
	if !ok {
		return 0, errors.New("last ID is not of type int32")
	}
	return lastID, nil
}

func (m *UserModel) IsAdmin(id int) (bool, error) {
	filter := bson.D{{Key: "id", Value: id}}

	var user bson.M

	err := m.Collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, ErrNoUserFound
		}
		return false, err
	}

	isAdmin, ok := user["isAdmin"].(bool)
	if !ok {
		return false, ErrNoUserFound
	}

	return isAdmin, nil
}
