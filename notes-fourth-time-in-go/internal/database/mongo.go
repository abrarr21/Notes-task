package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abrarr21/notes-in-golang/internal/config"
	"github.com/abrarr21/notes-in-golang/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	client *mongo.Client
	DB     *mongo.Database
	Users  *mongo.Collection
	Notes  *mongo.Collection
}

func ConnectDB(cfg *config.DatabaseConfig) (*Database, error) {
	// .SetServerSelectionTimeout() -> if MongoDB is unreachable, any query fails fast in 5 seconds instead of silently waiting 30 sec.
	// .SetMaxConnIdleTime() -> reclaim idle connections
	c, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoDB_URI).SetServerSelectionTimeout(5 * time.Second).SetMaxConnIdleTime(1 * time.Minute))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to Ping MongoDB: %w", err)

	}

	log.Println("connected to MongoDB ✅")

	db := c.Database(cfg.DBName)
	users := db.Collection("users")
	notes := db.Collection("notes")

	if err := models.EnsureIndexes(users); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return &Database{
		client: c,
		DB:     db,
		Users:  users,
		Notes:  notes,
	}, nil
}

func (d *Database) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := d.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to Disconnect from MongoDB: %w", err)
	}

	log.Println("Disconnected from MongoDB")
	return nil
}
