package database

import (
    "fmt"
    "log"
    "os"
    
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    
    "backend/internal/config"
    "backend/internal/models"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
    var dsn string
    
    databaseURL := os.Getenv("DATABASE_URL")
    if databaseURL != "" {
        dsn = databaseURL
        log.Println(" Using DATABASE_URL from environment")
    } else {
        dsn = fmt.Sprintf(
            "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
            cfg.DBHost,
            cfg.DBPort,
            cfg.DBUser,
            cfg.DBPass,
            cfg.DBName,
        )
        log.Println(" Using local database config")
    }
    
    var err error
    DB, err = gorm.Open(postgres.New(postgres.Config{
        DSN:                  dsn,
        PreferSimpleProtocol: true,
    }), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }
    
    log.Println(" Database connected successfully")
    return nil
}

func Migrate() error {
    log.Println("Running database migrations...")
    
    err := DB.AutoMigrate(
        &models.User{},
        &models.Topic{},
        &models.Post{},
        &models.Comment{},
        &models.Vote{}, 
    )
    
    if err != nil {
        return err
    }
    
    DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_vote_post_user ON votes(post_id, user_id, deleted_at)")
    
    return nil
}