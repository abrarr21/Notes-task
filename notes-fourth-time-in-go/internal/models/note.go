package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Database Model
type Note struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Title     string        `bson:"title"`
	Content   string        `bson:"content"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

// Incoming Request
type AddNoteRequest struct {
	Title   string `json:"title" validate:"required,min=4,max=20"`
	Content string `json:"content" validate:"required,min=10,max=25"`
}

// Outgoing API Response
type NoteResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateNoteRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}
