package main

import (
    "log"
    "gonote/internal/config"
    "gonote/internal/database"
    "gonote/internal/handler"
    "gonote/internal/middleware"
    "gonote/internal/repository"
    "gonote/internal/service"
    "gonote/internal/utils"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
    cfg := config.Load()
    
    db, err := database.NewPostgresDB(cfg)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    redisClient := database.NewRedisClient(cfg)
    
    userRepo := repository.NewUserRepository(db)
    noteRepo := repository.NewNoteRepository(db)
    cacheRepo := repository.NewCacheRepository(redisClient)
    
    jwtUtil := utils.NewJWTUtil(cfg.JWTSecret)
    authService := service.NewAuthService(userRepo, jwtUtil)
    aiService := service.NewAIService(cfg.GeminiAPIKey)
    noteService := service.NewNoteService(noteRepo, cacheRepo, aiService)
    
    authHandler := handler.NewAuthHandler(authService)
    noteHandler := handler.NewNoteHandler(noteService)
    
    router := gin.Default()
    
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))
    
    // Serve static files from frontend directory
    router.Static("/css", "../frontend/css")
    router.Static("/js", "../frontend/js")
    router.StaticFile("/", "../frontend/index.html")
    router.StaticFile("/login.html", "../frontend/login.html")
    router.StaticFile("/signup.html", "../frontend/signup.html")
    router.StaticFile("/dashboard.html", "../frontend/dashboard.html")
    
    // API routes
    v1 := router.Group("/api/v1")
    {
        auth := v1.Group("/auth")
        {
            auth.POST("/signup", authHandler.SignUp)
            auth.POST("/login", authHandler.Login)
        }
        
        notes := v1.Group("/notes")
        notes.Use(middleware.JWTAuthMiddleware(jwtUtil))
        {
            notes.POST("", noteHandler.CreateNote)
            notes.GET("", noteHandler.GetNotes)
            notes.GET("/:id", noteHandler.GetNote)
            notes.PUT("/:id", noteHandler.UpdateNote)
            notes.DELETE("/:id", noteHandler.DeleteNote)
            notes.POST("/:id/summary", noteHandler.GetSummary)
        }
    }
    
    log.Printf("Server starting on http://localhost:%s", cfg.ServerPort)
    router.Run(":" + cfg.ServerPort)
}