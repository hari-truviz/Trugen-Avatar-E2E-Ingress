package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	CustomTypes "infra-ingress/customtypes"
	Utilities "infra-ingress/utilities/infrastructure/runpod"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
	// Register websocket server
	http.HandleFunc("/ws", handleWebSocket)
	// Register server health endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Trugen Infra Ingress running."))
		w.WriteHeader(http.StatusOK)
	})
	fmt.Println("Web Socket Server is running on :8080/ws")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
