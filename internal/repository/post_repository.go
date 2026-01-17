package repository

import (
    "backend/internal/database"
    "backend/internal/models"
)

func GetPostsByTopicID(topicID uint, searchQuery string, sortBy string) ([]models.Post, error) {
    var posts []models.Post
    
    query := database.DB.Model(&models.Post{}).
        Preload("Creator").
        Where("topic_id = ?", topicID)
    
    if searchQuery != "" {
        searchTerm := "%" + searchQuery + "%"
        query = query.Where("title ILIKE ? OR content ILIKE ?", searchTerm, searchTerm)
    }
    
    if sortBy == "votes" {
        query = query.
            Joins("LEFT JOIN (SELECT post_id, SUM(vote_type) as vote_sum FROM votes WHERE deleted_at IS NULL GROUP BY post_id) AS vote_counts ON posts.id = vote_counts.post_id").
            Order("COALESCE(vote_counts.vote_sum, 0) DESC, posts.created_at DESC").
            Select("posts.*")
    } else {
        query = query.Order("created_at DESC")
    }

    result := query.Find(&posts)
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