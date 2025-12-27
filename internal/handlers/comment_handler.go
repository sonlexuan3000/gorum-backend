package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "backend/internal/middleware"
    "backend/internal/models"
    "backend/internal/repository"
)

type CreateCommentRequest struct {
    PostID   uint   `json:"post_id" binding:"required"`
    Content  string `json:"content" binding:"required,min=1,max=5000"`
    ParentID *uint  `json:"parent_id"`
}

type UpdateCommentRequest struct {
    Content string `json:"content" binding:"required,min=1,max=5000"`
}

func GetCommentsByPostID(c *gin.Context) {
    postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }
    
    _, err = repository.GetPostByID(uint(postID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
        return
    }
    
    comments, err := repository.GetCommentsByPostID(uint(postID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
        return
    }
    
    c.JSON(http.StatusOK, comments)
}

func GetCommentByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
        return
    }
    
    comment, err := repository.GetCommentByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
        return
    }
    
    c.JSON(http.StatusOK, comment)
}

func CreateComment(c *gin.Context) {
    var req CreateCommentRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    _, err := repository.GetPostByID(req.PostID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
        return
    }
    
    if req.ParentID != nil {
        _, err := repository.GetCommentByID(*req.ParentID)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Parent comment not found"})
            return
        }
    }
    
    userID := middleware.GetCurrentUserID(c)
    
    comment := models.Comment{
        PostID:    req.PostID,
        Content:   req.Content,
        CreatedBy: userID,
        ParentID:  req.ParentID,
    }
    
    if err := repository.CreateComment(&comment); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
        return
    }
    
    createdComment, _ := repository.GetCommentByID(comment.ID)
    
    c.JSON(http.StatusCreated, createdComment)
}


func UpdateComment(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
        return
    }
    
    var req UpdateCommentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    comment, err := repository.GetCommentByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
        return
    }
    
    userID := middleware.GetCurrentUserID(c)
    if comment.CreatedBy != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this comment"})
        return
    }
    
    comment.Content = req.Content
    
    if err := repository.UpdateComment(comment); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
        return
    }
    
    c.JSON(http.StatusOK, comment)
}

func DeleteComment(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
        return
    }
    
    comment, err := repository.GetCommentByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
        return
    }
    
    userID := middleware.GetCurrentUserID(c)
    if comment.CreatedBy != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this comment"})
        return
    }
    
    if err := repository.DeleteComment(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}