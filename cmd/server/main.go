package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/RyderDProgrammer/Inventory-Manager/internal/handlers"
	"github.com/RyderDProgrammer/Inventory-Manager/internal/middleware"
	"github.com/RyderDProgrammer/Inventory-Manager/internal/repository"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	repo := repository.NewInventoryRepository()
	h := handlers.NewHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.HealthCheck)
	mux.HandleFunc("GET /items", h.ListItems)
	mux.HandleFunc("GET /items/{id}", h.GetItem)
	mux.HandleFunc("POST /items", h.CreateItem)
	mux.HandleFunc("PUT /items/{id}", h.UpdateItem)
	mux.HandleFunc("DELETE /items/{id}", h.DeleteItem)

	slog.Info("server starting", "port", port, "env", env)
	if err := http.ListenAndServe(":"+port, middleware.Logging(mux)); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
