package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"real-time-chat-app/db"
)

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UserID    int64              `bson:"user_id" json:"user_id"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"password"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// Counter represents a sequence counter for generating user IDs
type Counter struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	Name   string             `bson:"name" json:"name"`
	Count  int64              `bson:"count" json:"count"`
	SeqNum int64              `bson:"seq_num" json:"seq_num"`
}

// UserService provides user-related operations

// UserService provides user-related operations
type UserService struct {
	userCollection *mongo.Collection
	counter        *mongo.Collection
}

// NewUserService creates a new UserService instance
func NewUserService() *UserService {
	return &UserService{
		userCollection: db.GetCollection("users"),
		counter:        db.GetCollection("counters"),
	}
}

// GetNextSequence returns the next sequence number for user IDs
func (s *UserService) GetNextSequence() (int64, error) {
	ctx := context.Background()
	filter := bson.M{"name": "user_id"}
	update := bson.M{
		"$inc": bson.M{"seq_num": int64(1)},
	}
	options := options.FindOneAndUpdate().SetUpsert(true)
	result := s.counter.FindOneAndUpdate(ctx, filter, update, options)
	if result.Err() != nil {
		return 0, result.Err()
	}

	var counter Counter
	if err := result.Decode(&counter); err != nil {
		return 0, err
	}

	return counter.SeqNum, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(username, password string) (*User, error) {
	ctx := context.Background()

	// Check if user already exists
	existingUser := User{}
	err := s.userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&existingUser)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	// Get next user ID
	userID, err := s.GetNextSequence()
	if err != nil {
		return nil, err
	}

	// Create new user
	user := User{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Username:  username,
		Password:  password, // In production, you should hash the password
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert user
	_, err = s.userCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(username string) (*User, error) {
	ctx := context.Background()
	user := User{}
	err := s.userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateLastLogin updates the user's last login time
func (s *UserService) UpdateLastLogin(username string) error {
	ctx := context.Background()
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err := s.userCollection.UpdateOne(ctx, bson.M{"username": username}, update)
	return err
}
