package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    
    "backend/internal/middleware"
    "backend/internal/models"
    "backend/internal/repository"
)

type VoteRequest struct {
    VoteType int `json:"vote_type" binding:"required,oneof=1 -1"` 
}

func VotePost(c *gin.Context) {
    postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }
    
    var req VoteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    _, err = repository.GetPostByID(uint(postID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
        return
    }
    
    userID := middleware.GetCurrentUserID(c)
    
    vote := models.Vote{
        PostID:   uint(postID),
        UserID:   userID,
        VoteType: req.VoteType,
    }
    
    if err := repository.UpsertVote(&vote); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to vote"})
        return
    }
    
    voteCount, _ := repository.GetVoteCount(uint(postID))
    userVote, _ := repository.GetUserVote(uint(postID), userID)
    
    c.JSON(http.StatusOK, gin.H{
        "vote_count": voteCount,
        "user_vote":  userVote,
    })
}

func UnvotePost(c *gin.Context) {
    postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }
    
    userID := middleware.GetCurrentUserID(c)
    
    if err := repository.DeleteVote(uint(postID), userID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove vote"})
        return
    }
    
    voteCount, _ := repository.GetVoteCount(uint(postID))
    
    c.JSON(http.StatusOK, gin.H{
        "vote_count": voteCount,
        "user_vote":  0,
    })
}