// This file is for topic-related database operations

package repository

import (
    "backend/internal/database"
    "backend/internal/models"
)

func GetAllTopics() ([]models.Topic, error) {
    var topics []models.Topic
    result := database.DB.Preload("Creator").Order("created_at DESC").Find(&topics)
    return topics, result.Error
}

func GetTopicByID(id uint) (*models.Topic, error) {
    var topic models.Topic
    result := database.DB.Preload("Creator").First(&topic, id)
    return &topic, result.Error
}

func CreateTopic(topic *models.Topic) error {
    return database.DB.Create(topic).Error
}

func UpdateTopic(topic *models.Topic) error {
    return database.DB.Save(topic).Error
}

func DeleteTopic(id uint) error {
    return database.DB.Delete(&models.Topic{}, id).Error
}