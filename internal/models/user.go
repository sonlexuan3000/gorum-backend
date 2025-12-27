package models

import (
    "time"
    "gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
    ID           uint           `gorm:"primarykey" json:"id"`
    Username     string         `gorm:"unique;not null" json:"username" validate:"required,min=3,max=50"`
    Email        string         `gorm:"unique;not null" json:"email" validate:"required,email"` 
    PasswordHash string         `gorm:"not null" json:"-"` 
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) SetPassword(password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    if err != nil {
        return err
    }
    u.PasswordHash = string(hashedPassword)
    return nil
}

func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
    return err == nil
}