package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *mongo.Database {
	// Connect to the MongoDB server.
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}

	// Access the "testdb" database.
	db := client.Database("testdb")

	// Read the setup script from the file.
	script, err := os.ReadFile("./testdata/setup.mongo")
	if err != nil {
		client.Disconnect(context.Background())
		t.Fatal(err)
	}

	// Execute the setup script to populate the database with test data.
	_, err = db.RunCommand(context.Background(), bson.M{"eval": script})
	if err != nil {
		client.Disconnect(context.Background())
		t.Fatal(err)
	}

	// Use t.Cleanup() to register a function to tear down the database after the tests are finished.
	t.Cleanup(func() {
		defer client.Disconnect(context.Background())
		script, err := os.ReadFile("./testdata/teardown.mongo")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.RunCommand(context.Background(), bson.M{"eval": script})

		if err != nil {
			t.Fatal(err)
		}
	})

	// Return the database.
	return db
}
