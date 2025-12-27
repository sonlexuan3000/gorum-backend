package repository

import (
    "backend/internal/database"
    "backend/internal/models"
)

func GetCommentsByPostID(postID uint) ([]models.Comment, error) {
    var comments []models.Comment
    
    result := database.DB.
        Preload("Creator").
        Preload("Replies.Creator").           
        Preload("Replies.Replies.Creator"). 
        Where("post_id = ? AND parent_id IS NULL", postID).
        Order("created_at DESC").
        Find(&comments)
    
    return comments, result.Error
}

func GetCommentByID(id uint) (*models.Comment, error) {
    var comment models.Comment
    result := database.DB.
        Preload("Creator").
        Preload("Post").
        First(&comment, id)
    return &comment, result.Error
}

func CreateComment(comment *models.Comment) error {
    return database.DB.Create(comment).Error
}

func UpdateComment(comment *models.Comment) error {
    return database.DB.Save(comment).Error
}

func DeleteComment(id uint) error {
    return database.DB.Delete(&models.Comment{}, id).Error
}

func GetRepliesByCommentID(commentID uint) ([]models.Comment, error) {
    var replies []models.Comment
    result := database.DB.
        Preload("Creator").
        Where("parent_id = ?", commentID).
        Order("created_at ASC").
        Find(&replies)
    return replies, result.Error
}