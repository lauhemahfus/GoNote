package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    ServerPort    string
    DBHost        string
    DBPort        string
    DBUser        string
    DBPassword    string
    DBName        string
    RedisHost     string
    RedisPort     string
    RedisPassword string
    JWTSecret     string
    GeminiAPIKey  string
}

func Load() *Config {
    godotenv.Load()
    
    return &Config{
        ServerPort:    getEnv("SERVER_PORT", "8080"),
        DBHost:        getEnv("DB_HOST", "localhost"),
        DBPort:        getEnv("DB_PORT", "5432"),
        DBUser:        getEnv("DB_USER", "postgres"),
        DBPassword:    getEnv("DB_PASSWORD", ""),
        DBName:        getEnv("DB_NAME", "gonote"),
        RedisHost:     getEnv("REDIS_HOST", "localhost"),
        RedisPort:     getEnv("REDIS_PORT", "6379"),
        RedisPassword: getEnv("REDIS_PASSWORD", ""),
        JWTSecret:     getEnv("JWT_SECRET", "secret-key"),
        GeminiAPIKey:  getEnv("GEMINI_API_KEY", ""),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}