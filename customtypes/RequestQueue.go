package CustomTypes

import (
	"container/list"
	"sync"

	"github.com/gorilla/websocket"
)

// Request
type Request struct {
	Conn *websocket.Conn
	Req  ConversationRequest
}

// Request queue
type RequestQueue struct {
	Mu       sync.Mutex
	Queue    *list.List
	IndexMap map[*websocket.Conn][]*list.Element
}

func (rq *RequestQueue) Enqueue(req Request) int {
	rq.Mu.Lock()
	defer rq.Mu.Unlock()
	elem := rq.Queue.PushBack(req)
	rq.IndexMap[req.Conn] = append(rq.IndexMap[req.Conn], elem)
	return rq.Queue.Len()
}

func (rq *RequestQueue) Dequeue() (Request, bool) {
	rq.Mu.Lock()
	defer rq.Mu.Unlock()
	front := rq.Queue.Front()
	if front == nil {
		return Request{}, false
	}
	req := front.Value.(Request)
	rq.Queue.Remove(front)
	// also clean up from indexMap
	elems := rq.IndexMap[req.Conn]
	for i, e := range elems {
		if e == front {
			rq.IndexMap[req.Conn] = append(elems[:i], elems[i+1:]...)
			break
		}
	}
	if len(rq.IndexMap[req.Conn]) == 0 {
		delete(rq.IndexMap, req.Conn)
	}
	return req, true
}

func (rq *RequestQueue) RemoveByConnection(conn *websocket.Conn) {
	rq.Mu.Lock()
	defer rq.Mu.Unlock()
	if elems, ok := rq.IndexMap[conn]; ok {
		for _, e := range elems {
			rq.Queue.Remove(e)
		}
		delete(rq.IndexMap, conn)
	}
}
