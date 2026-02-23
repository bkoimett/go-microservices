package server

import (
    // "encoding/json"
    "log"
    // "net/http"
    "sync"
	
    "github.com/gorilla/websocket"
    "github.com/bkoimett/go-microservices/crdt-sync-engine/internal/models"
)

// Client represents a connected WebSocket client
type Client struct {
    ID       string
    DocID    string
    Conn     *websocket.Conn
    Send     chan []byte
    Version  int64
}

// Hub manages all WebSocket connections
type Hub struct {
    clients    map[string]map[string]*Client // docID -> clientID -> Client
    register   chan *Client
    unregister chan *Client
    broadcast  chan *models.SyncMessage
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[string]map[string]*Client),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        broadcast:  make(chan *models.SyncMessage),
    }
}

// Run starts the hub's event loop
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            if _, ok := h.clients[client.DocID]; !ok {
                h.clients[client.DocID] = make(map[string]*Client)
            }
            h.clients[client.DocID][client.ID] = client
            h.mu.Unlock()
            log.Printf("Client %s registered for document %s", client.ID, client.DocID)

        case client := <-h.unregister:
            h.mu.Lock()
            if doc, ok := h.clients[client.DocID]; ok {
                if _, ok := doc[client.ID]; ok {
                    delete(doc, client.ID)
                    close(client.Send)
                    log.Printf("Client %s unregistered", client.ID)
                }
                if len(doc) == 0 {
                    delete(h.clients, client.DocID)
                }
            }
            h.mu.Unlock()

        case message := <-h.broadcast:
            h.mu.RLock()
            if clients, ok := h.clients[message.DocID]; ok {
                for _, client := range clients {
                    // Don't send back to the sender
                    if client.ID != message.ClientID {
                        select {
                        case client.Send <- message.Payload:
                        default:
                            close(client.Send)
                            delete(clients, client.ID)
                        }
                    }
                }
            }
            h.mu.RUnlock()
        }
    }
}