package routes

import (
	"github.com/abrarr21/test/internal/config"
	"github.com/abrarr21/test/internal/database"
	"github.com/abrarr21/test/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(db *database.Database, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	h := handlers.New(db, cfg)

	r.Get("/", h.CheckHealth)

	UserRoute(r, h)

	return r
}
