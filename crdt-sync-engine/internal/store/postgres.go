package store

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "github.com/bkoimett/go-microservices/crdt-sync-engine/internal/models"
)

type DocumentStore struct {
    db *sql.DB
}

func NewDocumentStore(connStr string) (*DocumentStore, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    // Create tables if they don't exist
    if err := createTables(db); err != nil {
        return nil, err
    }

    return &DocumentStore{db: db}, nil
}

func createTables(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS documents (
        id VARCHAR(36) PRIMARY KEY,
        content BYTEA NOT NULL,
        version BIGINT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );

    CREATE INDEX IF NOT EXISTS idx_documents_version ON documents(version);
    `

    _, err := db.Exec(query)
    return err
}

// SaveDocument persists a document or creates it if new
func (s *DocumentStore) SaveDocument(doc *models.Document) error {
    query := `
    INSERT INTO documents (id, content, version, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT (id) DO UPDATE SET
        content = EXCLUDED.content,
        version = EXCLUDED.version,
        updated_at = EXCLUDED.updated_at
    WHERE documents.version < EXCLUDED.version
    `

    result, err := s.db.Exec(query, doc.ID, doc.Content, doc.Version, doc.CreatedAt, doc.UpdatedAt)
    if err != nil {
        return fmt.Errorf("failed to save document: %w", err)
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        // Version was older, no update performed
        return nil
    }

    return nil
}

// GetDocument retrieves a document by ID
func (s *DocumentStore) GetDocument(id string) (*models.Document, error) {
    var doc models.Document
    
    query := `SELECT id, content, version, created_at, updated_at FROM documents WHERE id = $1`
    err := s.db.QueryRow(query, id).Scan(&doc.ID, &doc.Content, &doc.Version, &doc.CreatedAt, &doc.UpdatedAt)
    
    if err == sql.ErrNoRows {
        return nil, nil // Document doesn't exist
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get document: %w", err)
    }
    
    return &doc, nil
}

// GetDocumentsAfterVersion retrieves all documents updated after a certain version
// Useful for clients catching up after being offline [citation:7]
func (s *DocumentStore) GetDocumentsAfterVersion(version int64) ([]*models.Document, error) {
    query := `SELECT id, content, version, created_at, updated_at FROM documents WHERE version > $1`
    
    rows, err := s.db.Query(query, version)
    if err != nil {
        return nil, fmt.Errorf("failed to query documents: %w", err)
    }
    defer rows.Close()
    
    var documents []*models.Document
    for rows.Next() {
        var doc models.Document
        if err := rows.Scan(&doc.ID, &doc.Content, &doc.Version, &doc.CreatedAt, &doc.UpdatedAt); err != nil {
            return nil, err
        }
        documents = append(documents, &doc)
    }
    
    return documents, nil
}