package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID       bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string        `bson:"name" json:"name" validate:"required,alpha"`
	Email    string        `bson:"email" json:"email" validate:"required,email"`
	Password string        `bson:"password" json:"password" validate:"required,min=4"`
}

type UserRequest struct {
	Email    string `bson:"email" json:"email" validate:"required,email"`
	Password string `bson:"password" json:"password" validate:"required"`
}

type UserResponse struct {
	ID    bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string        `bson:"name" json:"name"`
	Email string        `bson:"email" json:"email"`
}

func EnsureIndexes(collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("unique_email_index"),
	}

	_, err := collection.Indexes().CreateOne(ctx, index)
	return err
}
