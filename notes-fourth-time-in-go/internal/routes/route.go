package routes

import (
	"github.com/abrarr21/notes-in-golang/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterAllRoutes(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", h.CheckHealth)

	UserRoutes(r, h)
	NoteRoutes(r, h)

	return r
}
