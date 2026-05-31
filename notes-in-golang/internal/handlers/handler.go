package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/abrarr21/notes/internal/database"
	"github.com/abrarr21/notes/internal/models"
	"github.com/abrarr21/notes/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Handler struct {
	DB *database.Database
}

func New(db *database.Database) *Handler {
	return &Handler{
		DB: db,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server running"))
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.Notes

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		log.Printf("failed to decode request body %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := utils.Validate.Struct(note); err != nil {
		log.Printf("validation failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	note.ID = bson.NewObjectID()
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	result, err := h.DB.Notes.InsertOne(ctx, note)
	if err != nil {
		log.Printf("failed to insert note: %v", err)
		http.Error(w, "failed to creat note", http.StatusInternalServerError)
		return
	}

	log.Println("Note inserted: ", result.InsertedID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"message": "note created successfully",
		"id":      result.InsertedID,
		"note":    note,
	}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *Handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := h.DB.Notes.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Printf("failed to fetch notes: %v", err)
		http.Error(w, "failed to fetch notes", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var notes []*models.Notes
	if err := cursor.All(ctx, &notes); err != nil {
		log.Printf("failed to decode notes: %v", err)
		http.Error(w, "failed to decode notes", http.StatusInternalServerError)
		return
	}

	if notes == nil {
		notes = []*models.Notes{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"message": "notes fetched successfully",
		"count":   len(notes),
		"notes":   notes,
	}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
