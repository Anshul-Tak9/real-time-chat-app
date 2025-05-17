package authentication

import (
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"real-time-chat-app/config"
	"real-time-chat-app/services"
)

// Login is the authenticator function for JWT
type LoginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type PasswordResetRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	Password2 string `json:"password2" binding:"required,eqfield=Password"`
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

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	// Return both user_id and username as a map
	userData := map[string]string{
		"user_id":  user.UserID,
		"username": user.Username,
	}

	return userData, nil
}

// ResetPassword handles password reset request
func ResetPassword(c *gin.Context) {
	var req PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userService := services.NewUserService()
	err := userService.ResetPassword(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password has been reset successfully",
	})
}

// Authorizator checks if the user is authorized
func Authorizator(data interface{}, c *gin.Context) bool {
	log.Println("Authorizing user:", data)
	userMap, ok := data.(map[string]string)
	if !ok {
		log.Println("Error: Invalid data type in Authorizator")
		return false
	}

	userId, ok := userMap["user_id"]
	if !ok {
		log.Println("Error: user_id not found in data")
		return false
	}

	username := userMap["username"]
	if username == "" {
		log.Println("Error: username not found in data")
		return false
	}

	// Store both values in context
	c.Set("user_id", userId)
	c.Set("username", username)
	return true
}

// Unauthorized handles unauthorized access
func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

func GetUserById(c *gin.Context) (interface{}, error) {
	var login map[string]string
	if err := c.ShouldBind(&login); err != nil {
		return nil, err
	}

	userService := services.NewUserService()
	userId := login[config.JWConfig.JWTIdentity]
	user, err := userService.GetUserById(userId)
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
			userData, ok := data.(map[string]string)
			if !ok {
				return jwt.MapClaims{"error": "invalid data type"}
			}
			return jwt.MapClaims{
				config.JWConfig.JWTIdentity: userData["user_id"],
				"username":                  userData["username"],
			}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return map[string]string{
				"user_id":  claims[config.JWConfig.JWTIdentity].(string),
				"username": claims["username"].(string),
			}
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
