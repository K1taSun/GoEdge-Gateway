package main

import (
	"log"
	"net/http"

	"github.com/k1tasun/GoEdge-Gateway/internal/config"
)

func main() {
	cfg := config.Load()

	log.Printf("Starting Storage Service on port %s", cfg.ServerPort)

	// TODO: Initialize Database connection
	// db := postgres.NewConnection(cfg.DatabaseURL)

	// TODO: Setup gRPC server or HTTP handlers

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

