package models

import (
    "time"
    "gorm.io/gorm"
)

type Topic struct {
    ID          uint           `gorm:"primarykey" json:"id"`
    Title       string         `gorm:"not null" json:"title" validate:"required,min=3,max=200"`
    Description string         `json:"description" validate:"max=1000"`
    CreatedBy   uint           `json:"created_by"`
    Creator     User           `gorm:"foreignKey:CreatedBy" json:"creator"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
    Posts       []Post         `gorm:"foreignKey:TopicID" json:"posts,omitempty"`
}