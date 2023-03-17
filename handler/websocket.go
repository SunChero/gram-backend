package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

var (
	connections = make(map[string]*clientConnection)
	mutex       = sync.Mutex{}
	// nc          *nats.Conn
)

type clientConnection struct {
	conn          *websocket.Conn
	clientID      string
	isAlive       bool
	subscriptions map[string]*nats.Subscription
}

func (h *Handler) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade to WebSocket connection
	upgrader := websocket.Upgrader{
		//CheckOrigin: checkOrigin,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Read client ID from query parameter
	clientID := r.URL.Query().Get("Id")
	if clientID == "" {
		log.Println("Missing clientId query parameter")
		conn.Close()
		return
	}

	// Lock the connections map
	mutex.Lock()

	// Check if client already has an active connection
	if connObj, ok := connections[clientID]; ok {
		// Close existing connection and replace with new one
		connObj.conn.Close()
	}

	// Create new client connection
	connections[clientID] = &clientConnection{
		conn:          conn,
		clientID:      clientID,
		isAlive:       true,
		subscriptions: make(map[string]*nats.Subscription),
	}

	// Subscribe to NATS subject
	sub, err := h.NC.Subscribe("socket."+clientID, func(msg *nats.Msg) {
		mutex.Lock()
		defer mutex.Unlock()

		if connObj, ok := connections[clientID]; ok && connObj.isAlive {
			if err := connObj.conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
				log.Printf("Failed to send message to client %s: %v", clientID, err)
			}
		}
	})
	if err != nil {
		log.Printf("Failed to subscribe to NATS subject: %v", err)
	} else {
		connections[clientID].subscriptions[sub.Subject] = sub
	}

	// Unlock the connections map
	mutex.Unlock()

	// Handle disconnections
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}

		// Lock the connections map
		mutex.Lock()

		for _, sub := range connections[clientID].subscriptions {
			sub.Unsubscribe()
		}
		conn.Close()
		delete(connections, clientID)

		// Unlock the connections map
		mutex.Unlock()
	}()
}
