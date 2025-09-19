package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	CustomTypes "infra-ingress/customtypes"
	Utilities "infra-ingress/utilities/infrastructure/runpod"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Handle avatar info route
func handleAvatarInfo(w http.ResponseWriter, r *http.Request) {
	avatarID := r.URL.Query().Get("id")
	userName := r.URL.Query().Get("userName")
	// Read landing avatars config json
	data, err := os.ReadFile("static/landing-page-avatars.json")
	if err != nil {
		fmt.Println("Unable to read file: %s", err)
	}
	var avatars = &CustomTypes.AllAvatars{}
	json.Unmarshal(data, avatars)
	// Filter the avatars
	for _, item := range avatars.Avatars {
		if item.AvatarID == avatarID {
			// Replace <USER_NAME> and <AVATAR_ID> tags
			item.PersonaPrompt = strings.ReplaceAll(item.PersonaPrompt, "<USER_NAME>", userName)
			item.PersonaPrompt = strings.ReplaceAll(item.PersonaPrompt, "<AVATAR_ID>", item.PersonaName)
			avatarJson, _ := json.Marshal(item)
			w.Write(avatarJson)
			w.WriteHeader(http.StatusOK)
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnection(conn *websocket.Conn, rq *CustomTypes.RequestQueue) {
	for {
		// Convert bytes[] to string to be unmarshelded and added to indexMap
		msgType, payload, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
		}
		defer func() {
			// Delete client request from the global map
			rq.RemoveByConnection(conn)
			conn.Close()
		}()
		convRequest := CustomTypes.ConversationRequest{}
		err = json.Unmarshal(payload, &convRequest)
		if err != nil {
			fmt.Println(err)
		}
		// Add conversation request to request queue
		currentPos := rq.Enqueue(CustomTypes.Request{
			Conn: conn,
			Req:  convRequest,
		})
		// Send acknowledgement message
		var ackMsg CustomTypes.Acknowledgment = CustomTypes.Acknowledgment{
			Message:         "new request accepted",
			CurrentPosition: currentPos,
			StatusCode:      201,
		}
		ack, err := json.Marshal(ackMsg)
		err = conn.WriteMessage(msgType, ack)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Handle message
	handleConnection(conn, rq)
}

// Instantiate Request Queue object
var rq = &CustomTypes.RequestQueue{
	Queue:    list.New(),
	IndexMap: make(map[*websocket.Conn][]*list.Element),
}

// Main entrypoint
func main() {
	err := godotenv.Load(".local.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	// Start a ticker to check the Runpod if it's ready
	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for range ticker.C {
			// Check if Infrastructure is ready to accept new requests
			if Utilities.IsRead() {
				// Pick the next item from the requestMap and post to runpod
				req, exist := rq.Dequeue()
				if exist {
					// Post the request to runpod
					Utilities.PostJob(req.Req)
				}
				fmt.Printf("Number of requests in the queue: %v\n", rq.Queue.Len())
			}
		}
	}()
	mux := http.NewServeMux()
	// Register Avatar Info route
	mux.HandleFunc("GET /avatars/", handleAvatarInfo)
	// Register websocket server
	mux.HandleFunc("/ws", handleWebSocket)
	// Register server health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Trugen Infra Ingress running."))
		w.WriteHeader(http.StatusOK)
	})
	// Register infra health
	mux.HandleFunc("/infra", func(w http.ResponseWriter, r *http.Request) {
		var statusCode int = http.StatusServiceUnavailable
		if Utilities.IsRead() {
			statusCode = http.StatusOK
		}
		w.WriteHeader(statusCode)
	})
	fmt.Println("Web Socket Server is running on :8080/ws")
	handler := cors.AllowAll().Handler(mux)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
