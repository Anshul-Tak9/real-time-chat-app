package router

import (
	"log"
	"net/http"
	"real-time-chat-app/authentication"
	"real-time-chat-app/chat"

	"github.com/gin-gonic/gin"
)

// Initialize initializes the router
func Initialize() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	// JWT middleware
	authMiddleware, _ := authentication.AuthMiddleware()

	// Public routes
	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/signup", authentication.SignUp)
	r.POST("/resetPassword", authentication.ResetPassword)
	r.POST("/getUserById", func(c *gin.Context) {
		user, err := authentication.GetUserById(c)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	})
	r.GET("/health", HealthCheckHandler)
	log.Println("Router initialized")
	// Protected routes
	api := r.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.POST("/createRoom", chat.CreateRoom)
		api.GET("/rooms", chat.ListRooms)
		api.GET("/rooms/:id/history", chat.RoomHistory)
	}

	// WebSocket (also protected)
	r.GET("/ws", authMiddleware.MiddlewareFunc(), chat.WSHandler)

	return r
}
