package repository

import (
    "backend/internal/database"
    "backend/internal/models"
)

func GetVoteCount(postID uint) (int, error) {
    var result struct {
        Total int
    }
    
    err := database.DB.Raw(`
        SELECT COALESCE(SUM(vote_type), 0) as total 
        FROM votes 
        WHERE post_id = ? AND deleted_at IS NULL
    `, postID).Scan(&result).Error
    
    return result.Total, err
}

func GetUserVote(postID, userID uint) (int, error) {
    var vote models.Vote
    err := database.DB.Where("post_id = ? AND user_id = ?", postID, userID).First(&vote).Error
    
    if err != nil {
        if err.Error() == "record not found" {
            return 0, nil 
        }
        return 0, err
    }
    
    return vote.VoteType, nil
}

func UpsertVote(vote *models.Vote) error {
    var existing models.Vote
    err := database.DB.Where("post_id = ? AND user_id = ?", vote.PostID, vote.UserID).First(&existing).Error
    
    if err != nil {
        return database.DB.Create(vote).Error
    }
    
    
    if existing.VoteType == vote.VoteType {
        return database.DB.Delete(&existing).Error
    }
    
    existing.VoteType = vote.VoteType
    return database.DB.Save(&existing).Error
}

func DeleteVote(postID, userID uint) error {
    return database.DB.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.Vote{}).Error
}