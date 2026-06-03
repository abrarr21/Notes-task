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
)

func (h *Handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Note hit"))
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
