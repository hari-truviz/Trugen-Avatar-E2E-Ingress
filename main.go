package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	CustomTypes "infra-ingress/customtypes"
	Utilities "infra-ingress/utilities/infrastructure/runpod"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConnection(conn *websocket.Conn, rq *RequestQueue) {
	for {
		// Convert bytes[] to string to be unmarshelded and added to indexMap
		_, payload, err := conn.ReadMessage()
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
		rq.Enqueue(Request{
			conn: conn,
			req:  convRequest,
		})
		// Send acknowledgement message
		var ackMsg string = "Conversation request added into the queue."
		err = conn.WriteMessage(1, []byte(ackMsg))
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

// Request
type Request struct {
	conn *websocket.Conn
	req  CustomTypes.ConversationRequest
}

// Request queue
type RequestQueue struct {
	mu       sync.Mutex
	queue    *list.List
	indexMap map[*websocket.Conn][]*list.Element
}

func (rq *RequestQueue) Enqueue(req Request) {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	elem := rq.queue.PushBack(req)
	rq.indexMap[req.conn] = append(rq.indexMap[req.conn], elem)
}

func (rq *RequestQueue) Dequeue() (Request, bool) {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	front := rq.queue.Front()
	if front == nil {
		return Request{}, false
	}
	req := front.Value.(Request)
	rq.queue.Remove(front)
	// also clean up from indexMap
	elems := rq.indexMap[req.conn]
	for i, e := range elems {
		if e == front {
			rq.indexMap[req.conn] = append(elems[:i], elems[i+1:]...)
			break
		}
	}
	if len(rq.indexMap[req.conn]) == 0 {
		delete(rq.indexMap, req.conn)
	}
	return req, true
}

func (rq *RequestQueue) RemoveByConnection(conn *websocket.Conn) {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	if elems, ok := rq.indexMap[conn]; ok {
		for _, e := range elems {
			rq.queue.Remove(e)
		}
		delete(rq.indexMap, conn)
	}
}

var rq = &RequestQueue{
	queue:    list.New(),
	indexMap: make(map[*websocket.Conn][]*list.Element),
}

// Main entrypoint
func main() {
	// Start a ticker to check the Runpod if it's ready
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			// Check if Infrastructure is ready to accept new requests
			if Utilities.IsRead() {
				// Pick the next item from the requestMap and post to runpod
				req, exist := rq.Dequeue()
				if exist {
					fmt.Println(req.req)
				}
				fmt.Printf("Number of requests in the queue: %v\n", rq.queue.Len())
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
