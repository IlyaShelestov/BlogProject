package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlockModelInterface interface {
	GetAll() ([]Block, error)
	Insert(title, description, imageURL1, imageURL2, imageURL3 string) error
	Update(id int, title, description, imageURL1, imageURL2, imageURL3 string) error
	Delete(id int) error
}

type Block struct {
	ID          int       `bson:"id"`
	CreatedTime time.Time `bson:"created_time"`
	ImageURL1   string    `bson:"image_url_1"`
	ImageURL2   string    `bson:"image_url_2"`
	ImageURL3   string    `bson:"image_url_3"`
	Description string    `bson:"description"`
	Title       string    `bson:"title"`
}

type BlockModel struct {
	Collection *mongo.Collection
}

func NewBlockModel(database *mongo.Database) *BlockModel {
	return &BlockModel{
		Collection: database.Collection("blocks"),
	}
}

func (m *BlockModel) GetAll() ([]Block, error) {
	var blocks []Block
	cursor, err := m.Collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &blocks); err != nil {
		return nil, err
	}
	return blocks, nil
}

func (m *BlockModel) Insert(title, description, imageURL1, imageURL2, imageURL3 string) error {
	block := bson.M{
		"created_time": time.Now(),
		"image_url_1":  imageURL1,
		"image_url_2":  imageURL2,
		"image_url_3":  imageURL3,
		"description":  description,
		"title":        title,
	}

	_, err := m.Collection.InsertOne(context.Background(), block)
	return err
}

func (m *BlockModel) Update(id int, title, description, imageURL1, imageURL2, imageURL3 string) error {
	filter := bson.D{{"id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"title", title},
			{"description", description},
			{"image_url_1", imageURL1},
			{"image_url_2", imageURL2},
			{"image_url_3", imageURL3},
		}},
	}

	_, err := m.Collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (m *BlockModel) Delete(id int) error {
	filter := bson.D{{"id", id}}
	_, err := m.Collection.DeleteOne(context.Background(), filter)
	return err
}
