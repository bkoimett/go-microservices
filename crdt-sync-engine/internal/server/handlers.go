package server

import (
    "encoding/json"
	"time"
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/google/uuid"
    "github.com/bkoimett/go-microservices/crdt-sync-engine/internal/models"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins for hackathon
    },
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

// HandleWebSocket upgrades HTTP to WebSocket
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    // Get document ID from URL query
    docID := r.URL.Query().Get("docId")
    if docID == "" {
        http.Error(w, "docId required", http.StatusBadRequest)
        return
    }

    // Generate or get client ID
    clientID := r.URL.Query().Get("clientId")
    if clientID == "" {
        clientID = uuid.New().String()
    }

    // Upgrade to WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade failed: %v", err)
        return
    }

    // Create client
    client := &Client{
        ID:      clientID,
        DocID:   docID,
        Conn:    conn,
        Send:    make(chan []byte, 256),
        Version: 0,
    }

    // Register client
    h.register <- client

    // Start goroutines for reading/writing
    go client.writePump()
    go client.readPump(h)
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
    defer func() {
        c.Conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.Send:
            if !ok {
                // Channel closed
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
                return
            }
        }
    }
}

// handleSync processes a sync request from a client
func (h *Hub) handleSync(client *Client, msg *models.SyncMessage) {
    // Get current document state from database
    doc, err := h.store.GetDocument(msg.DocID)
    if err != nil {
        log.Printf("Error getting document: %v", err)
        return
    }

    // If document doesn't exist, create it
    if doc == nil {
        doc = models.NewDocument()
        doc.ID = msg.DocID
        
        // If client sent initial content, use it
        if len(msg.Payload) > 0 {
            doc.Content = msg.Payload
        }
        
        if err := h.store.SaveDocument(doc); err != nil {
            log.Printf("Error saving new document: %v", err)
            return
        }
    }

    // Compare versions to determine what to send back
    response := &models.SyncMessage{
        Type:    "sync-response",
        DocID:   doc.ID,
        Version: doc.Version,
    }

    if msg.Version < doc.Version {
        // Client is behind, send full document
        response.Payload = doc.Content
        response.Type = "sync-full"
    } else if msg.Version > doc.Version {
        // Client is ahead (offline edits), we need to merge
        // For now, we'll just accept the client's version
        // In a real CRDT implementation, you'd merge here [citation:2]
        doc.Content = msg.Payload
        doc.Version = msg.Version
        doc.UpdatedAt = time.Now()
        h.store.SaveDocument(doc)
        
        response.Type = "sync-ack"
    }

    // Send response back to the requesting client
    responseJSON, _ := json.Marshal(response)
    client.Send <- responseJSON
}

func (c *Client) readPump(h *Hub) {
    defer func() {
        h.unregister <- c
        c.Conn.Close()
    }()

    c.Conn.SetReadLimit(512 * 1024) // 512KB max message size
    
    for {
        var msg models.SyncMessage
        err := c.Conn.ReadJSON(&msg)
        if err != nil {
            break
        }

        // Update client's version based on message
        if msg.Version > c.Version {
            c.Version = msg.Version
        }

        switch msg.Type {
        case "sync":
            h.handleSync(c, &msg)
        case "update":
            // Broadcast update to other clients
            h.broadcast <- &msg
            
            // Persist to database (in a real implementation, you'd batch these)
            go func(m models.SyncMessage) {
                doc, _ := h.store.GetDocument(m.DocID)
                if doc == nil {
                    doc = models.NewDocument()
                    doc.ID = m.DocID
                }
                
                // This is simplistic - real CRDT would merge here
                doc.Content = m.Payload
                doc.Version = m.Version
                doc.UpdatedAt = time.Now()
                
                h.store.SaveDocument(doc)
            }(msg)
        }
    }
}