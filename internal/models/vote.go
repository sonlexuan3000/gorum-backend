package models

import (
    "time"
    "gorm.io/gorm"
)

type Vote struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    PostID    uint           `gorm:"not null;index:idx_vote_post_user" json:"post_id"`
    Post      Post           `gorm:"foreignKey:PostID" json:"post,omitempty"`
    UserID    uint           `gorm:"not null;index:idx_vote_post_user" json:"user_id"`
    User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
    VoteType  int            `gorm:"not null" json:"vote_type"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
