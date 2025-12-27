package repository

import (
    "backend/internal/database"
    "backend/internal/models"
)

func GetPostsByTopicID(topicID uint) ([]models.Post, error) {
    var posts []models.Post
    result := database.DB.
        Preload("Creator").
        Where("topic_id = ?", topicID).
        Order("created_at DESC").
        Find(&posts)

    if result.Error != nil {
        return nil, result.Error
    }
    
    for i := range posts {
        voteCount, _ := GetVoteCount(posts[i].ID)
        posts[i].VoteCount = voteCount
    }
    
    return posts, nil
}

func GetPostByID(id uint) (*models.Post, error) {
    var post models.Post
    result := database.DB.
        Preload("Creator").
        Preload("Topic").
        First(&post, id)

    if result.Error != nil {
        return nil, result.Error
    }
    
    voteCount, _ := GetVoteCount(id)
    post.VoteCount = voteCount
    
    return &post, nil
}

func GetPostByIDWithUserVote(id, userID uint) (*models.Post, error) {
    post, err := GetPostByID(id)
    if err != nil {
        return nil, err
    }
    
    userVote, _ := GetUserVote(id, userID)
    post.UserVote = userVote
    
    return post, nil
}


func CreatePost(post *models.Post) error {
    return database.DB.Create(post).Error
}

func UpdatePost(post *models.Post) error {
    return database.DB.Save(post).Error
}

func DeletePost(id uint) error {
    return database.DB.Delete(&models.Post{}, id).Error
}

func CountPostsByTopicID(topicID uint) (int64, error) {
    var count int64
    result := database.DB.Model(&models.Post{}).Where("topic_id = ?", topicID).Count(&count)
    return count, result.Error
}