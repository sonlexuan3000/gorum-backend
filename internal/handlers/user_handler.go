package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    
    "backend/internal/middleware"
    "backend/internal/repository"
)

type UpdateProfileRequest struct {
    Bio       string `json:"bio" binding:"max=500"`
    AvatarURL string `json:"avatar_url" binding:"omitempty,url"`
}

func GetCurrentUser(c *gin.Context) {
    userID := middleware.GetCurrentUserID(c)
    
    user, err := repository.GetUserWithStats(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    c.JSON(http.StatusOK, user)
}

func GetUserProfile(c *gin.Context) {
    identifier := c.Param("identifier")
    
    var userID uint
    var err error

    if id, parseErr := strconv.ParseUint(identifier, 10, 32); parseErr == nil {
        userID = uint(id)
    } else {
        user, err := repository.GetUserByUsername(identifier)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        userID = user.ID
    }
    
    user, err := repository.GetUserWithStats(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    c.JSON(http.StatusOK, user)
}

func UpdateProfile(c *gin.Context) {
    userID := middleware.GetCurrentUserID(c)
    
    var req UpdateProfileRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, err := repository.GetUserByID(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    user.Bio = req.Bio
    if req.AvatarURL != "" {
        user.AvatarURL = req.AvatarURL
    }
    
    if err := repository.UpdateUser(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
        return
    }
    
    c.JSON(http.StatusOK, user)
}


func GetUserPosts(c *gin.Context) {
    identifier := c.Param("identifier")
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
    
    var userID uint
    if id, err := strconv.ParseUint(identifier, 10, 32); err == nil {
        userID = uint(id)
    } else {
        user, err := repository.GetUserByUsername(identifier)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        userID = user.ID
    }
    
    posts, err := repository.GetPostsByUserID(userID, limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
        return
    }
    
    c.JSON(http.StatusOK, posts)
}

func GetUserComments(c *gin.Context) {
    identifier := c.Param("identifier")
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
    
    var userID uint
    if id, err := strconv.ParseUint(identifier, 10, 32); err == nil {
        userID = uint(id)
    } else {
        user, err := repository.GetUserByUsername(identifier)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        userID = user.ID
    }
    
    comments, err := repository.GetCommentsByUserID(userID, limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
        return
    }
    
    c.JSON(http.StatusOK, comments)
}