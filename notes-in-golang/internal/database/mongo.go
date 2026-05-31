package database

import (
	"context"
	"log"
	"time"

	"github.com/abrarr21/notes/internal/config"
	"github.com/abrarr21/notes/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	client *mongo.Client
	DB     *mongo.Database
	Notes  *mongo.Collection
}

func ConnectDB(cfg *config.DatabaseConfig) *Database {
	c, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoDB_URI))
	if err != nil {
		log.Fatal("failed to connect to MongoDB: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.Ping(ctx, nil); err != nil {
		log.Fatal("failed to reach MongoDB, check if it is running: ", err)
	}

	log.Println("Conncted to MongoDB")

	db := c.Database(cfg.DBName)
	notes := db.Collection("notes")

	models.EnsureNoteIndexes(notes)

	return &Database{
		client: c,
		DB:     db,
		Notes:  notes,
	}
}

func (d *Database) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := d.client.Disconnect(ctx); err != nil {
		log.Printf("failed disconnecting from MongoDB: %v", err)
	}

	log.Println("Disconnected from MongoDB")
}
