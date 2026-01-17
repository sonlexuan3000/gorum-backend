package repository

import (
    "backend/internal/database"
    "backend/internal/models"
)

func FindUserByEmail(email string) (*models.User, error) {
    var user models.User
    result := database.DB.Where("email = ?", email).First(&user)
    return &user, result.Error
}

func FindUserByUsername(username string) (*models.User, error) {
    var user models.User
    result := database.DB.Where("username = ?", username).First(&user)
    return &user, result.Error
}

func CreateUser(user *models.User) error {
    return database.DB.Create(user).Error
}

func GetUserByID(id uint) (*models.User, error) {
    var user models.User
    result := database.DB.First(&user, id)
    return &user, result.Error
}

func CheckEmailExists(email string) (bool, error) {
    var count int64
    err := database.DB.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
    return count > 0, err
}

func CheckUsernameExists(username string) (bool, error) {
    var count int64
    err := database.DB.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
    return count > 0, err
}

func GetUserByUsername(username string) (*models.User, error) {
    var user models.User
    result := database.DB.Where("username = ?", username).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}


func GetUserWithStats(id uint) (*models.User, error) {
    user, err := GetUserByID(id)
    if err != nil {
        return nil, err
    }
    
    var postsCount int64
    database.DB.Model(&models.Post{}).Where("created_by = ?", id).Count(&postsCount)
    user.PostsCount = int(postsCount)
    
    var commentsCount int64
    database.DB.Model(&models.Comment{}).Where("created_by = ?", id).Count(&commentsCount)
    user.CommentsCount = int(commentsCount)
    
    return user, nil
}

func UpdateUser(user *models.User) error {
    return database.DB.Save(user).Error
}


func GetPostsByUserID(userID uint, limit, offset int) ([]models.Post, error) {
    var posts []models.Post
    query := database.DB.Model(&models.Post{}).
        Preload("Topic").
        Preload("Creator").
        Where("created_by = ?", userID).
        Order("created_at DESC")
    
    if limit > 0 {
        query = query.Limit(limit).Offset(offset)
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

func GetCommentsByUserID(userID uint, limit, offset int) ([]models.Comment, error) {
    var comments []models.Comment
    query := database.DB.Model(&models.Comment{}).
        Preload("Post").
        Preload("Post.Topic").
        Preload("Creator").
        Where("created_by = ?", userID).
        Order("created_at DESC")
    
    if limit > 0 {
        query = query.Limit(limit).Offset(offset)
    }
    
    result := query.Find(&comments)
    return comments, result.Error
}