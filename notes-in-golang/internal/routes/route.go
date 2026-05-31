package routes

import (
	"github.com/abrarr21/notes/internal/database"
	"github.com/abrarr21/notes/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(db *database.Database) *chi.Mux {
	r := chi.NewRouter()

	h := handlers.New(db)
	r.Use(middleware.Logger)

	r.Get("/", h.Get)
	r.Post("/api/notes/", h.CreateNote)
	r.Get("/api/notes", h.GetAllNotes)

	return r
}
