package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    Port      string  
    DBHost    string  
    DBPort    string  
    DBUser    string
    DBPass    string
    DBName    string
    JWTSecret string  
}

func LoadConfig() (*Config, error) {
    
    godotenv.Load()

    return &Config{
        Port:      getEnv("PORT", "8080"),         
        DBHost:    getEnv("DB_HOST", "localhost"),
        DBPort:    getEnv("DB_PORT", "5432"),
        DBUser:    getEnv("DB_USER", "lexuanson"),
        DBPass:    getEnv("DB_PASSWORD", ""),
        DBName:    getEnv("DB_NAME", "forum_db"),
        JWTSecret: getEnv("JWT_SECRET", "secret"),
    }, nil
}

func getEnv(key, defaultVal string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultVal
}