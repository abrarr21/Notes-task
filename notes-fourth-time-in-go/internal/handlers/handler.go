package handlers

import (
	"net/http"

	"github.com/abrarr21/notes-in-golang/internal/database"
)

type Handler struct {
	DB *database.Database
}

func New(db *database.Database) *Handler {
	return &Handler{
		DB: db,
	}
}

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server running perfectly ✅"))
}
