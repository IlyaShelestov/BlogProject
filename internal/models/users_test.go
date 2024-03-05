package models

import (
	"blog_project/internal/assert"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestUserModel_Exists(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
			require.NoError(t, err)
			defer client.Disconnect(context.Background())
			collection := client.Database("testdb").Collection("users")
			m := UserModel{collection}
			exists, err := m.Exists(tt.userID)
			assert.Equal(t, tt.want, exists)
			assert.NilError(t, err)
		})
	}
}
