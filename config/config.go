package config

import (
	"time"
)

// Config holds the application configuration
var JWConfig = struct {
	// JWT Settings
	JWTSecret   string
	JWTTimeout  time.Duration
	JWTRealm    string
	JWTIdentity string

	// Server Settings
	Port string
}{
	JWTSecret:   "hkoT7ApfFDsNlur1v/d7cYRKzDVJpX4ugJtqBurYtLc=", // Secure JWT secret
	JWTTimeout:  time.Hour,
	JWTRealm:    "chat",
	JWTIdentity: "username",
}

var MongoDBConfig = struct {
	// MongoDB Settings
	MongoURI     string
	DatabaseName string
	MaxPoolSize  int
	MinPoolSize  int
	Timeout      int // in seconds
}{

	MongoURI:     "mongodb://localhost:27017",
	DatabaseName: "chat_app",
	MaxPoolSize:  100,
	MinPoolSize:  10,
	Timeout:      10,
}
