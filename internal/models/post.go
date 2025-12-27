package models

import (
    "time"
    "gorm.io/gorm"
)

type Post struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    TopicID   uint           `json:"topic_id" validate:"required"`
    Topic     Topic          `gorm:"foreignKey:TopicID" json:"topic,omitempty"`
    Title     string         `gorm:"not null" json:"title" validate:"required,min=3,max=200"`
    Content   string         `gorm:"type:text" json:"content" validate:"required"`
    CreatedBy uint           `json:"created_by"`
    Creator   User           `gorm:"foreignKey:CreatedBy" json:"creator"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    Comments  []Comment      `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	VoteCount int            `gorm:"-" json:"vote_count"`
    UserVote  int            `gorm:"-" json:"user_vote"`
}