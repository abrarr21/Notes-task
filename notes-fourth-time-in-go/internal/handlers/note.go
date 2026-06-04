package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/abrarr21/notes-in-golang/internal/models"
	"github.com/abrarr21/notes-in-golang/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (h *Handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := h.DB.Notes.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Printf("failed to fetch notes: %v", err)
		http.Error(w, "failed to fetch notes", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var notes []models.Note

	if err := cursor.All(ctx, &notes); err != nil {
		log.Printf("failed to decode notes: %v", err)
		http.Error(w, "failed to decode notes", http.StatusInternalServerError)
		return
	}

	var response []models.NoteResponse

	for _, note := range notes {
		response = append(response, models.NoteResponse{
			ID:        note.ID.Hex(),
			Title:     note.Title,
			Content:   note.Content,
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		})
	}

	if err := utils.ResponseJSON(w, http.StatusOK, "notes fetched successfully", response); err != nil {
		log.Println("failed to encode response")
	}
}

func (h *Handler) CreaetNotes(w http.ResponseWriter, r *http.Request) {
	var input models.Note

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println("error decoding note input")
		utils.ResponseJSON(w, http.StatusBadRequest, "invalid input", nil)
		return
	}

	if err := utils.Validator.Struct(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//your app might be deployed across servers or time zones
	now := time.Now().UTC()

	note := models.Note{
		ID:        bson.NewObjectID(),
		Title:     strings.TrimSpace(input.Title),
		Content:   strings.TrimSpace(input.Content),
		CreatedAt: now,
		UpdatedAt: now,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	result, err := h.DB.Notes.InsertOne(ctx, note)
	if err != nil {
		log.Println("failed to create note")
		http.Error(w, "failed to create note", http.StatusInternalServerError)
		return
	}

	log.Println("Note create: ", result.InsertedID)

	response := models.NoteResponse{
		ID:        note.ID.Hex(),
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}

	if err := utils.ResponseJSON(w, http.StatusCreated, "note created successfully", response); err != nil {
		log.Println("failed to encode respone")
	}
}
