package config

import (
	"os"
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
	JWTSecret:   os.Getenv("JWT_SECRET"),
	JWTTimeout:  time.Hour,
	JWTRealm:    "chat",
	JWTIdentity: "user_id",
}

var MongoDBConfig = struct {
	// MongoDB Settings
	MongoURI     string
	DatabaseName string
	MaxPoolSize  int
	MinPoolSize  int
	Timeout      int // in seconds
}{

	MongoURI:     os.Getenv("MONGODB_URI"),
	DatabaseName: os.Getenv("MONGODB_DATABASE"),
	MaxPoolSize:  100,
	MinPoolSize:  10,
	Timeout:      10,
}
