package main

import (
    "log"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    
    "backend/internal/config"
    "backend/internal/database"
	"backend/internal/routes"
)

func main() {
    log.Println("Starting Forum Backend Server...")
    
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    if err := database.Connect(cfg); err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    if err := database.Migrate(); err != nil {
        log.Fatal("Failed to migrate database:", err)
    }
    
    r := gin.Default()
    
    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{
            "http://localhost:5173",           // Change this to your frontend URL
            os.Getenv("FRONTEND_URL"),         // Deployment frontend URL
        },
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))
    
	routes.SetupRoutes(r)

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "message": "Server is running!",
        })
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = cfg.Port
    }
    
    log.Printf("Server is running on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}