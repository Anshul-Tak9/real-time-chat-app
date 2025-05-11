# Real-Time Chat Application

A real-time chat application built with Go using Gin framework, MongoDB, and WebSocket for real-time communication.

## Features

- User Authentication (Signup/Login)
- Real-time Chat Rooms
- WebSocket-based Communication
- JWT-based Authentication
- MongoDB Integration

## Tech Stack

- Backend: Go (Gin Framework)
- Database: MongoDB
- Real-time: WebSocket
- Authentication: JWT

## Setup

### Using Docker (Recommended)

1. Prerequisites:
   - Docker
   - Docker Compose

2. Clone the repository:
   ```bash
   git clone <repository-url>
   cd real-time-chat-app
   ```

3. For production:
   ```bash
   docker-compose up --build
   ```

4. For development with hot-reloading:
   ```bash
   docker-compose up --build
   ```

The application will be available at `http://localhost:4000`

#### Debug Mode Features
- Hot-reloading with Air
- Automatic code reloading on changes
- Persistent MongoDB data
- Environment variables support

### Without Docker

1. Prerequisites:
   - Go 1.18+
   - MongoDB
   - Git

2. Clone the repository:
   ```bash
   git clone <repository-url>
   cd real-time-chat-app
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run MongoDB:
   ```bash
   mongod
   ```

5. Start the application:
   ```bash
   go run main.go
   ```

## API Endpoints

### Authentication
- POST `/signup` - Create a new user
- POST `/login` - Authenticate user

### Chat Rooms
- POST `/api/rooms` - Create a new chat room
- GET `/api/rooms` - List all chat rooms
- GET `/api/rooms/:id/history` - Get room chat history

### WebSocket
- GET `/ws` - WebSocket endpoint for real-time communication

## Database Schema

### Users Collection
- `_id`: MongoDB ObjectID
- `user_id`: Unique user identifier
- `username`: User's username
- `password`: User's password (Note: In production, passwords should be hashed)
- `created_at`: Timestamp of user creation
- `updated_at`: Timestamp of last update

### Counters Collection
- Used for generating unique user IDs

## Security Notes

- Passwords are currently stored in plain text (for development only)
- In production, passwords should be hashed using bcrypt or similar
- JWT tokens are used for authentication
- WebSocket connections are protected by JWT middleware

## Contributing

This is a personal project. Contributions are welcome but not required.

## License

MIT License
