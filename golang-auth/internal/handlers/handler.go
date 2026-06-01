package handlers

import (
	"net/http"

	"github.com/abrarr21/test/internal/config"
	"github.com/abrarr21/test/internal/database"
)

type Handlers struct {
	DB  *database.Database
	Cfg *config.Config
}

func New(db *database.Database, cfg *config.Config) *Handlers {
	return &Handlers{
		DB:  db,
		Cfg: cfg,
	}
}

func (h *Handlers) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server running perfectly"))
}
