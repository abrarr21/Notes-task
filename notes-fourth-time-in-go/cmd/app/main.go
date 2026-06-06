package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abrarr21/notes-in-golang/internal/config"
	"github.com/abrarr21/notes-in-golang/internal/database"
	"github.com/abrarr21/notes-in-golang/internal/handlers"
	"github.com/abrarr21/notes-in-golang/internal/routes"
)

func main() {
	cfg := config.Load()
	db, err := database.ConnectDB(&cfg.Database)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer func() {
		if err := db.Disconnect(); err != nil {
			log.Println("warning: could not Disconnect cleanly: ", err)
		}
	}()

	h := handlers.New(db, cfg)

	router := routes.RegisterAllRoutes(h)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("server listening on port: ", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println("failed to start server: ", err)
		}
	}()

	sig := <-quit
	log.Println("shutdown signal received: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("failed to close server gracefully, forcing shutdown: %v", err)
	}

	log.Println("server closed gracefully")
}
