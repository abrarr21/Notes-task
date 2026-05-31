package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abrarr21/notes/internal/config"
	"github.com/abrarr21/notes/internal/database"
	"github.com/abrarr21/notes/internal/routes"
)

func main() {
	cfg := config.Load()
	db := database.ConnectDB(&cfg.Database)
	defer db.Disconnect()

	r := routes.RegisterRoutes(db)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server listening on port:", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("server failed to start:", err)
		}
	}()

	sig := <-quit
	log.Println("shutdown signal received:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed, forcing exit: %v", err)
	}

	log.Println("Server stopped gracefully")
}
