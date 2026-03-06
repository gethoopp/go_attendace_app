# HR Attendance Application

A Go-based HR attendance management system with real-time RFID tracking, push notifications, and AI chatbot capabilities.

## Features

- **User Authentication**: Register, login, and logout with JWT tokens
- **Attendance Tracking**: Real-time RFID-based attendance via WebSocket
- **Check-in/Check-out**: Employee check-in and check-out functionality
- **Presence Management**: View and manage attendance records
- **Push Notifications**: Firebase Cloud Messaging (FCM) for notifications
- **AI Chatbot**: Ollama-powered chatbot for employee assistance
- **Image Recognition**: OpenCV integration for face/image recognition
- **Message Queue**: RabbitMQ for async task processing

## Tech Stack

- **Language**: Go 1.24
- **Framework**: Gin
- **Database**: MySQL 8.0
- **Authentication**: JWT
- **Push Notifications**: Firebase FCM
- **AI**: Ollama
- **Image Processing**: GoCV (OpenCV)
- **Message Queue**: RabbitMQ
- **Container**: Docker & Docker Compose

## Project Structure

```
go_attendace_app/
├── main.go                    # Application entry point
├── docker-compose.yaml        # Docker Compose configuration
├── dockerfile                 # Docker build file
├── install.sh                 # Installation script
├── .env                       # Environment variables
├── database/
│   └── get_database.go       # Database connection
├── modules/
│   ├── users.go              # User model
│   ├── attendance_request.go # Attendance request model
│   ├── device_token.go       # Device token model
│   ├── message_user.go       # Message model
│   └── chat_request.go       # Chat request model
├── services/
│   ├── auth.go               # Authentication service
│   ├── presence.go           # Presence tracking service
│   ├── services.go           # General services
│   └── jwt.go                # JWT service
├── middleware/
│   ├── jwt_middleware.go     # JWT middleware
│   └── fcm_middleware.go     # FCM middleware
├── push_notification/
│   ├── fcm.go                # Firebase Cloud Messaging
│   └── publisher_rabbitmq.go # RabbitMQ publisher
├── chat/
│   └── chatbot.go            # Ollama chatbot integration
└── image_recognition/
    └── imageGoCv.go          # OpenCV image recognition
```

## Prerequisites

- Go 1.24+
- Docker & Docker Compose
- MySQL 8.0 (or use Docker)
- Ollama (or use Docker)

## Installation

### Using Docker (Recommended)

```bash
# Build and run all services
docker-compose up -d

# Or build with custom configuration
docker-compose up --build
```

### Manual Installation

```bash
# Install dependencies
go mod download

# Run the application
go run main.go
```

## Environment Variables

Create a `.env` file with the following variables:

```env
PORT=8080
DB_USER=root
DB_PASS=your_password
DB_HOST=localhost
DB_NAME=go_attendance_app
JWT_SECRET=your_jwt_secret
FCM_API_KEY=your_fcm_api_key
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Health check |
| GET | `/ws/input` | WebSocket for RFID input |
| GET | `/api/data` | Get user data |
| GET | `/api/presence` | Get presence records |
| GET | `/api/totalwork` | Get total work data |
| POST | `/api/register` | Register new user |
| POST | `/api/login` | User login |
| POST | `/api/logout` | User logout |
| POST | `/api/getByDate` | Get presence by date |
| POST | `/api/saveToken` | Save device token |
| POST | `/api/sendNotif` | Send push notification |
| POST | `/api/chat` | AI chatbot endpoint |
| POST | `/api/checkIn` | Employee check-in |
| PUT | `/api/checkOut` | Employee check-out |

## Docker Services

- **mysql**: MySQL 8.0 database
- **ollama**: AI chatbot service
- **go-backend**: Main application

## License

MIT
