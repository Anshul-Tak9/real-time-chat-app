package authentication

import (
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"real-time-chat-app/config"
	"real-time-chat-app/services"
)

// Login is the authenticator function for JWT
type LoginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Login is the authenticator function for JWT
func Login(c *gin.Context) (interface{}, error) {
	var login LoginRequest
	if err := c.ShouldBind(&login); err != nil {
		return nil, err
	}

	userService := services.NewUserService()
	user, err := userService.GetUserByUsername(login.Username)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	// In production, you should compare hashed passwords
	if user.Password != login.Password {
		return nil, jwt.ErrFailedAuthentication
	}

	// Update last login time
	userService.UpdateLastLogin(login.Username)

	return user.Username, nil
}

// Authorizator checks if the user is authorized
func Authorizator(data interface{}, c *gin.Context) bool {
	log.Println("Authorizing user:", data)
	username, ok := data.(string)
	if !ok {
		log.Println("Error: Invalid data type in Authorizator")
		return false
	}

	// Verify the user exists in the database
	userService := services.NewUserService()
	user, err := userService.GetUserByUsername(username)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return false
	}
	if user == nil {
		log.Printf("Error: User not found in database for username: %s", username)
		return false
	}
	log.Println("User exists in the database:", userService)

	return true
}

// Unauthorized handles unauthorized access
func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

func GetUserByUsername(c *gin.Context) (interface{}, error) {
	var login map[string]string
	if err := c.ShouldBind(&login); err != nil {
		return nil, err
	}

	userService := services.NewUserService()
	user, err := userService.GetUserByUsername(login[config.JWConfig.JWTIdentity])
	if err != nil {
		return nil, err
	}

	return user, nil
}

// SignUp creates a new user
func SignUp(c *gin.Context) {
	var signUp LoginRequest
	if err := c.ShouldBind(&signUp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userService := services.NewUserService()
	user, err := userService.CreateUser(signUp.Username, signUp.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"user_id": user.UserID,
	})
}

func AuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       config.JWConfig.JWTRealm,
		Key:         []byte(config.JWConfig.JWTSecret),
		Timeout:     config.JWConfig.JWTTimeout,
		IdentityKey: config.JWConfig.JWTIdentity,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			return jwt.MapClaims{config.JWConfig.JWTIdentity: data.(string)}
		},
		Authenticator: Login,
		Authorizator:  Authorizator,
		Unauthorized:  Unauthorized,
	})
	if err != nil {
		return nil, err
	}
	return authMiddleware, nil
}
