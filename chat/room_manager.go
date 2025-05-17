package chat

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// RoomManager manages WebSocket connections for multiple rooms
type RoomManager struct {
	rooms     map[string]*Room
	roomMutex sync.RWMutex
}

// Room represents a single chat room
type Room struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	messages   []Message
	roomMutex  sync.RWMutex
}

// NewRoomManager creates a new room manager
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

// NewRoom creates a new room
func NewRoom() *Room {
	return &Room{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// JoinRoom adds a client to a room
func (rm *RoomManager) JoinRoom(roomId string, conn *websocket.Conn) {
	log.Printf("[ROOM] Client joining room %s", roomId)

	rm.roomMutex.Lock()
	room, exists := rm.rooms[roomId]
	if !exists {
		room = NewRoom()
		rm.rooms[roomId] = room
		go room.start()
		log.Printf("[ROOM] Created new room %s", roomId)
	}
	rm.roomMutex.Unlock()

	room.register <- conn
	log.Printf("[ROOM] Client registered in room %s", roomId)
}

// LeaveRoom removes a client from a room
func (rm *RoomManager) LeaveRoom(roomId string, conn *websocket.Conn) {
	rm.roomMutex.Lock()
	room, exists := rm.rooms[roomId]
	rm.roomMutex.Unlock()

	if exists {
		room.unregister <- conn
	}
}

// BroadcastMessage sends a message to all clients in a room
func (rm *RoomManager) BroadcastMessage(roomId string, message []byte) {
	rm.roomMutex.RLock()
	room, exists := rm.rooms[roomId]
	rm.roomMutex.RUnlock()

	if exists {
		log.Printf("[ROOM] Broadcasting message to room %s: %s", roomId, string(message))
		room.broadcast <- message
	} else {
		log.Printf("[ROOM] Room %s not found for broadcast", roomId)
	}
}

// start runs the room's message handling goroutine
func (r *Room) start() {
	log.Printf("[ROOM] Starting room goroutine")

	for {
		select {
		case conn := <-r.register:
			r.roomMutex.Lock()
			r.clients[conn] = true
			log.Printf("[ROOM] Client registered. Total clients: %d", len(r.clients))
			r.roomMutex.Unlock()
		case conn := <-r.unregister:
			r.roomMutex.Lock()
			if _, ok := r.clients[conn]; ok {
				delete(r.clients, conn)
				log.Printf("[ROOM] Client unregistered. Total clients: %d", len(r.clients))
				conn.Close()
			}
			r.roomMutex.Unlock()
		case message := <-r.broadcast:
			r.roomMutex.RLock()
			log.Printf("[ROOM] Broadcasting message to %d clients", len(r.clients))
			for conn := range r.clients {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("[ERROR] Failed to write message to client: %v", err)
					r.unregister <- conn
				} else {
					log.Printf("[ROOM] Message sent successfully to client")
				}
			}
			r.roomMutex.RUnlock()
		}
	}
}
