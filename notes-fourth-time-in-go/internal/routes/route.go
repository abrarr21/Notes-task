package routes

import (
	"github.com/abrarr21/notes-in-golang/internal/config"
	"github.com/abrarr21/notes-in-golang/internal/database"
	"github.com/abrarr21/notes-in-golang/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterAllRoutes(db *database.Database, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	h := handlers.New(db, cfg)

	r.Use(middleware.Logger)
	r.Get("/", h.CheckHealth)

	UserRoutes(r, h)
	NoteRoutes(r, h)

	return r
}
