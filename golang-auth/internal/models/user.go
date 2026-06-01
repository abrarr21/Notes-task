package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Address struct {
	City    string `bson:"city" json:"city"`
	Pincode int    `bson:"pincode" json:"pincode"`
}

type User struct {
	ID       bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string        `bson:"name" json:"name" validate:"required,min=3,max=20"`
	Email    string        `bson:"email" json:"email" validate:"required,email"`
	Password string        `bson:"password" json:"password" validate:"required"`
	Address  []Address     `bson:"address" json:"address"`
}

type UserResponse struct {
	ID      bson.ObjectID `json:"id"`
	Name    string        `json:"name"`
	Email   string        `json:"email"`
	Address []Address     `json:"address"`
}

func CreateIndexes(collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("unique_email_index"),
	}

	_, err := collection.Indexes().CreateOne(ctx, index)
	return err
}
