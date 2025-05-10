package authentication

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

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
	claims, ok := data.(jwt.MapClaims)
	if !ok {
		return false
	}
	
	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return false
	}
	
	// TODO: Implement actual authorization logic
	return true
}

// Unauthorized handles unauthorized access
func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
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
