package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    
    "backend/internal/middleware"
    "backend/internal/repository"
)

func GetNotifications(c *gin.Context) {
    userID := middleware.GetCurrentUserID(c)
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
    
    notifications, err := repository.GetUserNotifications(userID, limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
        return
    }
    
    c.JSON(http.StatusOK, notifications)
}

func GetUnreadCount(c *gin.Context) {
    userID := middleware.GetCurrentUserID(c)
    
    count, err := repository.GetUnreadNotificationsCount(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch count"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

func MarkAsRead(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
        return
    }
    
    if err := repository.MarkNotificationAsRead(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as read"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Marked as read"})
}

func MarkAllAsRead(c *gin.Context) {
    userID := middleware.GetCurrentUserID(c)
    
    if err := repository.MarkAllNotificationsAsRead(userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark all as read"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "All marked as read"})
}