package routes

import (
	"github.com/abrarr21/notes-in-golang/internal/handlers"
	"github.com/abrarr21/notes-in-golang/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

func NoteRoutes(r chi.Router, h *handlers.Handler) {
	r.Route("/api/notes", func(r chi.Router) {
		r.Get("/", h.GetAllNotes)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.RequireAuth(h.Cfg.JWT.JWT_SECRET))
			r.Post("/", h.CreaetNotes)
		})

	})
}
