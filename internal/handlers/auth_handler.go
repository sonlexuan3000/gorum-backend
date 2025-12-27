package handlers

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    
    "backend/internal/models"
    "backend/internal/repository"
    "backend/internal/utils"
)

type SignupRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6,max=100"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
    Token string       `json:"token"`
    User  models.User  `json:"user"`
}

func Signup(c *gin.Context) {
    var req SignupRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    req.Email = strings.ToLower(strings.TrimSpace(req.Email))
    req.Username = strings.TrimSpace(req.Username)
    
    emailExists, err := repository.CheckEmailExists(req.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    if emailExists {
        c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
        return
    }
    
    usernameExists, err := repository.CheckUsernameExists(req.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    if usernameExists {
        c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
        return
    }
    
    user := models.User{
        Username: req.Username,
        Email:    req.Email,
    }
    
    if err := user.SetPassword(req.Password); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
        return
    }
    
    if err := repository.CreateUser(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }
    
    token, err := utils.GenerateToken(user.ID, user.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }
    
    c.JSON(http.StatusCreated, AuthResponse{
        Token: token,
        User:  user,
    })
}

func Login(c *gin.Context) {
    var req LoginRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    req.Email = strings.ToLower(strings.TrimSpace(req.Email))
    
    user, err := repository.FindUserByEmail(req.Email)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
    
    if !user.CheckPassword(req.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }
    
    token, err := utils.GenerateToken(user.ID, user.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }
    
    c.JSON(http.StatusOK, AuthResponse{
        Token: token,
        User:  *user,
    })
}

func GetMe(c *gin.Context) {
    userID, _ := c.Get("user_id")
    
    user, err := repository.GetUserByID(userID.(uint))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    c.JSON(http.StatusOK, user)
}