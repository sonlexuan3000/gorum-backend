package repository

import (
    "backend/internal/database"
    "backend/internal/models"
)

func CreateNotification(notification *models.Notification) error {
    if notification.UserID == notification.ActorID {
        return nil
    }
    return database.DB.Create(notification).Error
}

func GetUserNotifications(userID uint, limit, offset int) ([]models.Notification, error) {
    var notifications []models.Notification
    query := database.DB.Model(&models.Notification{}).
        Preload("Actor").
        Preload("Post").
        Preload("Post.Topic").
        Preload("Comment").
        Where("user_id = ?", userID).
        Order("created_at DESC")
    
    if limit > 0 {
        query = query.Limit(limit).Offset(offset)
    }
    
    result := query.Find(&notifications)
    return notifications, result.Error
}

func GetUnreadNotificationsCount(userID uint) (int64, error) {
    var count int64
    result := database.DB.Model(&models.Notification{}).
        Where("user_id = ? AND is_read = ?", userID, false).
        Count(&count)
    return count, result.Error
}

func MarkNotificationAsRead(id uint) error {
    return database.DB.Model(&models.Notification{}).
        Where("id = ?", id).
        Update("is_read", true).Error
}

func MarkAllNotificationsAsRead(userID uint) error {
    return database.DB.Model(&models.Notification{}).
        Where("user_id = ? AND is_read = ?", userID, false).
        Update("is_read", true).Error
}

func DeleteNotification(id uint) error {
    return database.DB.Delete(&models.Notification{}, id).Error
}