package routes

import (
	"github.com/abrarr21/test/internal/handlers"
	"github.com/abrarr21/test/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

func UserRoute(r chi.Router, h *handlers.Handlers) {
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/", h.CreateUser)
		r.Post("/login", h.Login)

		// single with auth middlewares
		r.With(middlewares.RequireAuth(h.Cfg.JWT.Secret)).Get("/me", h.Me)

		// group routes
		r.Group(func(r chi.Router) {
			r.Use(middlewares.RequireAuth(h.Cfg.JWT.Secret))

			r.Get("/me", h.Me)
		})
	})
}
