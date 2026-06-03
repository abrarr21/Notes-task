package routes

import (
	"github.com/abrarr21/notes-in-golang/internal/database"
	"github.com/abrarr21/notes-in-golang/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterAllRoutes(db *database.Database) *chi.Mux {
	r := chi.NewRouter()

	h := handlers.New(db)

	r.Get("/", h.CheckHealth)

	return r
}
