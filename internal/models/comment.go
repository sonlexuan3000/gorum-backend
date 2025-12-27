package models

import (
    "time"
    "gorm.io/gorm"
)

type Comment struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    PostID    uint           `json:"post_id" validate:"required"`
    Post      Post           `gorm:"foreignKey:PostID" json:"post,omitempty"`
    Content   string         `gorm:"type:text" json:"content" validate:"required,min=1,max=5000"`
    CreatedBy uint           `json:"created_by"`
    Creator   User           `gorm:"foreignKey:CreatedBy" json:"creator"`
    ParentID  *uint          `json:"parent_id"` 
    Parent    *Comment       `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
    Replies   []Comment      `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}