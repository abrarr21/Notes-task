package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/abrarr21/notes-in-golang/internal/models"
	"github.com/abrarr21/notes-in-golang/internal/utils"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

	response := make([]models.NoteResponse, 0)

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
	var input models.AddNoteRequest

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println("error decoding note input")
		utils.ResponseJSON(w, http.StatusBadRequest, "invalid input", nil)
		return
	}

	if err := utils.Validate(input); err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, "validation failed", err)
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

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	noteID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, "invalid note id", nil)
		return
	}

	var input models.UpdateNoteRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, "invalid input", nil)
		return
	}
	if errors := utils.Validate(input); errors != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, "validation failed", errors)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// fetch existing note
	var existing models.Note
	if err := h.DB.Notes.FindOne(ctx, bson.M{"_id": noteID}).Decode(&existing); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.ResponseJSON(w, http.StatusNotFound, "note not found", nil)
			return
		}
		log.Printf("failed to fetch note %s: %v", id, err)
		http.Error(w, "failed to fetch note", http.StatusInternalServerError)
		return
	}

	// reusable response builder
	toResponse := func(n models.Note) models.NoteResponse {
		return models.NoteResponse{
			ID:        n.ID.Hex(),
			Title:     n.Title,
			Content:   n.Content,
			CreatedAt: n.CreatedAt,
			UpdatedAt: n.UpdatedAt,
		}
	}

	// diff check — extend this map as fields grow
	if !hasChanges(map[*string]string{
		input.Title:   existing.Title,
		input.Content: existing.Content,
	}) {
		// utils.ResponseJSON(w, http.StatusOK, "no changes detected", toResponse(existing))
		utils.ResponseJSON(w, http.StatusOK, "no changes detected", nil)
		return
	}

	update, err := buildUpdateDoc(input)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	result, err := h.DB.Notes.UpdateOne(ctx, bson.M{"_id": noteID}, bson.M{"$set": update})
	if err != nil {
		log.Printf("failed to update note %s: %v", id, err)
		http.Error(w, "failed to update note", http.StatusInternalServerError)
		return
	}
	if result.MatchedCount == 0 {
		utils.ResponseJSON(w, http.StatusNotFound, "note not found", nil)
		return
	}

	// fetch and return the updated note
	var updated models.Note
	if err := h.DB.Notes.FindOne(ctx, bson.M{"_id": noteID}).Decode(&updated); err != nil {
		log.Printf("failed to fetch updated note %s: %v", id, err)
		http.Error(w, "failed to fetch updated note", http.StatusInternalServerError)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusOK, "note updated successfully", toResponse(updated)); err != nil {
		log.Printf("failed to encode response for note %s: %v", id, err)
	}
}

func (h *Handler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	noteID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, "invalid note id", nil)
		return
	}

	var note models.Note

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if err := h.DB.Notes.FindOne(ctx, bson.M{"_id": noteID}).Decode(&note); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.ResponseJSON(w, http.StatusNotFound, "note not found", nil)
			return
		}
		log.Printf("failed to fetch note %s: %v", id, err)
		http.Error(w, "failed to fetch note", http.StatusInternalServerError)
		return
	}

	response := models.NoteResponse{
		ID:        note.ID.Hex(),
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}

	if err := utils.ResponseJSON(w, http.StatusOK, "note fetched successfully", response); err != nil {
		log.Println("error encodring response")
	}

}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	noteID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, "invalid note id", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	res, err := h.DB.Notes.DeleteOne(ctx, bson.M{"_id": noteID})
	if err != nil {
		log.Printf("failed to delete note %s: %v", id, err)
		utils.ResponseJSON(w, http.StatusInternalServerError, "failed to delete note", nil)
		return
	}

	if res.DeletedCount == 0 {
		utils.ResponseJSON(w, http.StatusNotFound, "note not found", nil)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusOK, "note deleted", nil); err != nil {
		log.Println("error encoding response")
	}
}
