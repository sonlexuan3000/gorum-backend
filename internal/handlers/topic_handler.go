package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    
    "backend/internal/middleware"
    "backend/internal/models"
    "backend/internal/repository"
)

type CreateTopicRequest struct {
    Title       string `json:"title" binding:"required,min=3,max=200"`
    Description string `json:"description" binding:"max=1000"`
}

type UpdateTopicRequest struct {
    Title       string `json:"title" binding:"required,min=3,max=200"`
    Description string `json:"description" binding:"max=1000"`
}

func GetAllTopics(c *gin.Context) {
    topics, err := repository.GetAllTopics()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch topics"})
        return
    }
    
    c.JSON(http.StatusOK, topics)
}

func GetTopicByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
        return
    }
    
    topic, err := repository.GetTopicByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
        return
    }
    
    c.JSON(http.StatusOK, topic)
}

func CreateTopic(c *gin.Context) {
    var req CreateTopicRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    userID := middleware.GetCurrentUserID(c)
    
    topic := models.Topic{
        Title:       req.Title,
        Description: req.Description,
        CreatedBy:   userID,
    }
    
    if err := repository.CreateTopic(&topic); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create topic"})
        return
    }
    
    createdTopic, _ := repository.GetTopicByID(topic.ID)
    
    c.JSON(http.StatusCreated, createdTopic)
}


func UpdateTopic(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
        return
    }
    
    var req UpdateTopicRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    topic, err := repository.GetTopicByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
        return
    }
    
    userID := middleware.GetCurrentUserID(c)
    if topic.CreatedBy != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this topic"})
        return
    }
    
    topic.Title = req.Title
    topic.Description = req.Description
    
    if err := repository.UpdateTopic(topic); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update topic"})
        return
    }
    
    c.JSON(http.StatusOK, topic)
}

func DeleteTopic(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
        return
    }
    
    topic, err := repository.GetTopicByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
        return
    }
    
    userID := middleware.GetCurrentUserID(c)
    if topic.CreatedBy != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this topic"})
        return
    }
    
    if err := repository.DeleteTopic(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete topic"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Topic deleted successfully"})
}