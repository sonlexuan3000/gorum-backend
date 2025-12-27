package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    
    "backend/internal/middleware"
    "backend/internal/models"
    "backend/internal/repository"
)

type CreatePostRequest struct {
    TopicID uint   `json:"topic_id" binding:"required"`
    Title   string `json:"title" binding:"required,min=3,max=200"`
    Content string `json:"content" binding:"required,min=1"`
}

type UpdatePostRequest struct {
    Title   string `json:"title" binding:"required,min=3,max=200"`
    Content string `json:"content" binding:"required,min=1"`
}

func GetPostsByTopicID(c *gin.Context) {
    topicID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
        return
    }
    
    _, err = repository.GetTopicByID(uint(topicID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
        return
    }
    
    posts, err := repository.GetPostsByTopicID(uint(topicID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
        return
    }
    
    c.JSON(http.StatusOK, posts)
}

func GetPostByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }
    
    post, err := repository.GetPostByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
        return
    }
    
    c.JSON(http.StatusOK, post)
}

func CreatePost(c *gin.Context) {
    var req CreatePostRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    _, err := repository.GetTopicByID(req.TopicID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
        return
    }
    
    userID := middleware.GetCurrentUserID(c)
    
    post := models.Post{
        TopicID:   req.TopicID,
        Title:     req.Title,
        Content:   req.Content,
        CreatedBy: userID,
    }
    
    if err := repository.CreatePost(&post); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
        return
    }
    
    createdPost, _ := repository.GetPostByID(post.ID)
    
    c.JSON(http.StatusCreated, createdPost)
}

func UpdatePost(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }
    
    var req UpdatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    post, err := repository.GetPostByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
        return
    }
    
    userID := middleware.GetCurrentUserID(c)
    if post.CreatedBy != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this post"})
        return
    }
    
    post.Title = req.Title
    post.Content = req.Content
    
    if err := repository.UpdatePost(post); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
        return
    }
    
    c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }
    
    post, err := repository.GetPostByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
        return
    }
    
    userID := middleware.GetCurrentUserID(c)
    if post.CreatedBy != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this post"})
        return
    }
    
    if err := repository.DeletePost(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}