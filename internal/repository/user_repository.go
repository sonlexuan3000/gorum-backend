
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