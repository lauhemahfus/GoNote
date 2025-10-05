# GoNote - Personal Note Taking Application

A modern, fast, and AI-powered note-taking application built with Go, Gin, PostgreSQL, Redis, and Google Gemini AI. Features a beautiful blue pastel UI theme and follows clean architecture principles.

## Features

- **User Authentication** - Secure JWT-based authentication system
- **CRUD Operations** - Create, Read, Update, and Delete notes effortlessly
- **Redis Caching** - Lightning-fast note retrieval with intelligent caching
- **AI-Powered Summaries** - Generate instant summaries using Google Gemini 2.0 Flash
- **Lazy Loading** - Efficient pagination for handling large note collections


## Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/lauhemahfus/gonote.git
cd gonote
```

### 2. Set Up PostgreSQL Database

```bash
# Create database
createdb gonote
# Run migrations
psql -d gonote -f backend/migrations/001_init.sql
```

### 3. Set Up Redis

```bash
# Start Redis server
redis-server
```

### 4. Configure Environment Variables

```bash
# Copy the example environment file
cp .env.example .env
# Edit .env with your credentials
```

### 5. Install Go Dependencies

```bash
cd backend
# Download dependencies
go mod download
# Tidy up modules
go mod tidy
```

**Or install manually:**

```bash
go get github.com/gin-gonic/gin@v1.9.1
go get github.com/gin-contrib/cors@v1.4.0
go get github.com/go-redis/redis/v8@v8.11.5
go get github.com/golang-jwt/jwt/v4@v4.5.0
go get github.com/joho/godotenv@v1.5.1
go get github.com/lib/pq@v1.10.9
go get golang.org/x/crypto@v0.14.0
```

### 6. Run the Application

```bash
go run cmd/server/main.go
```

You should see:
```
Server starting on http://localhost:8080
```

### 7. Access the Application

Open your browser and visit:

```
http://localhost:8080
```


## Screenshots

<img width="1333" height="856" alt="image" src="https://github.com/user-attachments/assets/58a105a1-3969-4ba8-a869-fe40d59dbee4" />
<img width="1439" height="859" alt="image" src="https://github.com/user-attachments/assets/5a158e66-e874-480a-8931-722d412e0501" />

