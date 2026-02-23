package models

import (
    "time"
    "github.com/google/uuid"
)

// Document represents a collaborative document
type Document struct {
    ID        string    `json:"id"`
    Content   []byte    `json:"content"` // Binary CRDT state
    Version   int64     `json:"version"`  // Monotonic version for sync
    UpdatedAt time.Time `json:"updated_at"`
    CreatedAt time.Time `json:"created_at"`
}

// SyncMessage represents messages between client and server
type SyncMessage struct {
    Type     string `json:"type"`               // "sync", "update", "ack"
    DocID    string `json:"docId"`
    Version  int64  `json:"version"`             // Client's current version
    Payload  []byte `json:"payload,omitempty"`   // CRDT update
    ClientID string `json:"clientId"`            // Unique client identifier
}

// NewDocument creates a new document with default values
func NewDocument() *Document {
    return &Document{
        ID:        uuid.New().String(),
        Content:   []byte{},  // Empty initial state
        Version:   0,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}