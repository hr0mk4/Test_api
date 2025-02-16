package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hr0mk4/test_api/internal/database"
	"github.com/hr0mk4/test_api/internal/models"
)

var jwtSecret = []byte("secret")

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func AuthHandler(c *gin.Context) {
	dbInstance := database.GetDB()
	if dbInstance == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Database connection error"})
		return
	}
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Invalid data"})
		return
	}

	var user models.User
	err := dbInstance.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		user = models.User{
			Username: req.Username,
			Password: req.Password,
			Coins:    1000,
		}
		if err := dbInstance.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "Create user error"})
			return
		}
	} else if user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "Wrong password"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Ошибка создания токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
