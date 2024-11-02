package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Message string `json:"message"`
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients    = make(map[*websocket.Conn]bool)
	clientsMux sync.Mutex
)

func broadcast(message Message) {
	clientsMux.Lock()
	defer clientsMux.Unlock()

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return
	}

	log.Printf("Broadcasting message: %s", string(jsonMessage))
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			log.Printf("Error sending message to client: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}
	defer ws.Close()

	// クライアント接続を記録
	clientsMux.Lock()
	clients[ws] = true
	clientCount := len(clients)
	clientsMux.Unlock()

	log.Printf("New client connected. Total clients: %d", clientCount)

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Client disconnected: %v", err)
			clientsMux.Lock()
			delete(clients, ws)
			log.Printf("Remaining clients: %d", len(clients))
			clientsMux.Unlock()
			break
		}

		log.Printf("Received message: %s", string(data))

		var msg Message
		err = json.Unmarshal(data, &msg)
		if err != nil {
			log.Printf("Error parsing JSON: %v. Using raw message", err)
			msg = Message{
				Message: string(data),
			}
		}

		broadcast(msg)
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	log.Println("Server starting at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
