package main

import (
    "log"
    "net/http"
    "os"
    "github.com/joho/godotenv"

    "github.com/bkoimett/go-microservices/crdt-sync-engine/internal/server"
    "github.com/bkoimett/go-microservices/crdt-sync-engine/internal/store"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    // Database connection
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        connStr = "postgres://postgres:password@localhost:5432/crdt_sync?sslmode=disable"
    }

    // Initialize store
    documentStore, err := store.NewDocumentStore(connStr)
    if err != nil {
        log.Fatalf("Failed to initialize store: %v", err)
    }
    defer documentStore.Close()

    // Create hub
    hub := server.NewHub()
    hub.SetStore(documentStore) // You'll need to add this method

    // Start hub in its own goroutine
    go hub.Run()

    // Set up routes
    http.HandleFunc("/ws", hub.HandleWebSocket)
    
    // Health check endpoint
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // Get port from environment or use default
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on :%s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}