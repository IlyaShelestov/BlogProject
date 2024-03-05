package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlockModelInterface interface {
	GetAll() ([]Block, error)
	Insert(title, description, imageURL1, imageURL2, imageURL3 string) error
	Update(id int, title, description, imageURL1, imageURL2, imageURL3 string) error
	Delete(id int) error
}

type Block struct {
	ID          int       `bson:"id"`
	Created     time.Time `bson:"created_time"`
	ImageURL1   string    `bson:"image_url_1"`
	ImageURL2   string    `bson:"image_url_2"`
	ImageURL3   string    `bson:"image_url_3"`
	Description string    `bson:"description"`
	Title       string    `bson:"title"`
}

type BlockModel struct {
	Collection *mongo.Collection
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
	ctx := context.Background()

	var result bson.M
	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
	err := m.Collection.FindOne(ctx, bson.D{}, opts).Decode(&result)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	fmt.Println(result)

	var nextID int = 0
	if result != nil {
		currentMaxID := int(result["id"].(int32))
		nextID = currentMaxID + 1
	}

	fmt.Println(nextID)

	block := bson.M{
		"id":           nextID,
		"created_time": time.Now(),
		"image_url_1":  imageURL1,
		"image_url_2":  imageURL2,
		"image_url_3":  imageURL3,
		"description":  description,
		"title":        title,
	}

	_, err = m.Collection.InsertOne(ctx, block)
	return err
}

func (m *BlockModel) Update(id int, title, description, imageURL1, imageURL2, imageURL3 string) error {
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "title", Value: title},
			{Key: "description", Value: description},
			{Key: "image_url_1", Value: imageURL1},
			{Key: "image_url_2", Value: imageURL2},
			{Key: "image_url_3", Value: imageURL3},
		}},
	}

	_, err := m.Collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (m *BlockModel) Delete(id int) error {
	filter := bson.D{{Key: "id", Value: id}}
	_, err := m.Collection.DeleteOne(context.Background(), filter)
	return err
}
