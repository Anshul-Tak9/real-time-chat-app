package router

import (
	"real-time-chat-app/authentication"
	"real-time-chat-app/chat"
	"real-time-chat-app/config"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// Setup returns a Gin engine with all middleware and routes registered.
func Initialize() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	// JWT middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       config.JWConfig.JWTRealm,
		Key:         []byte(config.JWConfig.JWTSecret),
		Timeout:     config.JWConfig.JWTTimeout,
		IdentityKey: config.JWConfig.JWTIdentity,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			return jwt.MapClaims{config.JWConfig.JWTIdentity: data.(string)}
		},
		Authenticator: authentication.Login,
		Authorizator:  authentication.Authorizator,
		Unauthorized:  authentication.Unauthorized,
	})
	if err != nil {
		panic("JWT Error:" + err.Error())
	}

	// Public routes
	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/signup", authentication.SignUp)

	// Protected routes
	api := r.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.POST("/rooms", chat.CreateRoom)
		api.GET("/rooms", chat.ListRooms)
		api.GET("/rooms/:id/history", chat.RoomHistory)
	}

	// WebSocket (also protected)
	r.GET("/ws", authMiddleware.MiddlewareFunc(), chat.WSHandler)

	return r
}
