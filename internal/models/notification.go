package models

import (
    "time"
    "gorm.io/gorm"
)

type NotificationType string

const (
    NotificationVote    NotificationType = "vote"
    NotificationComment NotificationType = "comment"
    NotificationReply   NotificationType = "reply"
)

type Notification struct {
    ID        uint             `gorm:"primarykey" json:"id"`
    UserID    uint             `gorm:"not null;index" json:"user_id"`           
    User      User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
    
    ActorID   uint             `gorm:"not null" json:"actor_id"`                
    Actor     User             `gorm:"foreignKey:ActorID" json:"actor"`
    
    Type      NotificationType `gorm:"type:varchar(20);not null" json:"type"`
    PostID    *uint            `gorm:"index" json:"post_id,omitempty"`          
    Post      *Post            `gorm:"foreignKey:PostID" json:"post,omitempty"`
    CommentID *uint            `gorm:"index" json:"comment_id,omitempty"`       
    Comment   *Comment         `gorm:"foreignKey:CommentID" json:"comment,omitempty"`
    
    IsRead    bool             `gorm:"default:false;index" json:"is_read"`
    CreatedAt time.Time        `json:"created_at"`
    UpdatedAt time.Time        `json:"updated_at"`
    DeletedAt gorm.DeletedAt   `gorm:"index" json:"-"`
}