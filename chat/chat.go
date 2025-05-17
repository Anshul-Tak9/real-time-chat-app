package chat

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"real-time-chat-app/db"
)

var roomManager = NewRoomManager()
var messageCollection *mongo.Collection

func init() {
	// Initialize MongoDB connection
	if err := db.InitializeDB(); err != nil {
		log.Fatal("Failed to initialize MongoDB:", err)
	}
	messageCollection = db.GetCollection("messages")
}

// WSHandler handles websocket connections

// CreateRoom creates a new room
func CreateRoom(c *gin.Context) {
	// TODO: Implement room creation logic
	c.JSON(http.StatusOK, gin.H{"message": "room created"})
}

// ListRooms lists all rooms
func ListRooms(c *gin.Context) {
	// TODO: Implement room listing logic
	c.JSON(http.StatusOK, gin.H{"rooms": []string{"room1", "room2"}})
}

// RoomHistory gets room history
func RoomHistory(c *gin.Context) {
	// TODO: Implement room history logic
	c.JSON(http.StatusOK, gin.H{"messages": []string{"message1", "message2"}})
}

func WSHandler(c *gin.Context) {
	roomId := c.Param("roomId")
	if roomId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room ID is required"})
		return
	}

	// Get user ID from context
	value, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
		return
	}
	userId := value.(string)

	usernamevalue, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found"})
		return
	}
	username := usernamevalue.(string)

	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[WS] Connection upgraded for room %s", roomId)
	defer roomManager.LeaveRoom(roomId, ws)

	// Join the room
	roomManager.JoinRoom(roomId, ws)

	// Get room history from MongoDB
	ctx := context.Background()
	var messages []Message
	cursor, err := messageCollection.Find(ctx, bson.M{"room_id": roomId})
	if err != nil {
		log.Printf("[ERROR] Failed to get room history: %v", err)
	} else {
		cursor.All(ctx, &messages)
		// Send history to new client
		for _, msg := range messages {
			if err := ws.WriteJSON(msg); err != nil {
				log.Printf("[ERROR] Failed to send history message: %v", err)
				break
			}
		}
	}
	// Get room history
	roomManager.roomMutex.RLock()
	room, exists := roomManager.rooms[roomId]
	roomManager.roomMutex.RUnlock()

	if exists {
		// Send history to new client
		for _, msg := range room.messages {
			if err := ws.WriteJSON(msg); err != nil {
				log.Printf("[ERROR] Failed to send history message: %v", err)
				break
			}
		}
	}

	// Send welcome message
	welcomeMsg := fmt.Sprintf("[SYSTEM] %s has entered the %s", username, roomId)
	roomManager.BroadcastMessage(roomId, []byte(welcomeMsg))

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			log.Printf("[WS] Connection closed for room %s", roomId)
			break
		}

		// Store message in room
		roomManager.roomMutex.Lock()
		if room, exists := roomManager.rooms[roomId]; exists {
			room.messages = append(room.messages, Message{
				RoomID:    roomId,
				UserID:    userId,
				Username:  username,
				Content:   string(message),
				CreatedAt: time.Now(),
			})
		}
		roomManager.roomMutex.Unlock()

		// Store message in MongoDB
		msg := Message{
			RoomID:    roomId,
			UserID:    userId,
			Username:  username,
			Content:   string(message),
			CreatedAt: time.Now(),
		}
		_, err = messageCollection.InsertOne(ctx, msg)
		if err != nil {
			log.Printf("[ERROR] Failed to store message in MongoDB: %v", err)
		}

		// Add prefix to incoming message for logging
		incomingMsg := fmt.Sprintf("[IN] Room %s: %s", roomId, string(message))
		log.Println("Received message:", incomingMsg)

		// Broadcast message to all clients in the room
		outgoingMsg := fmt.Sprintf("[OUT] Room %s: %s", roomId, string(message))
		log.Printf("[WS] Broadcasting message: %s", outgoingMsg)
		roomManager.BroadcastMessage(roomId, []byte(outgoingMsg))
	}
}
