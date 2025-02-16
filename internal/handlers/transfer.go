package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hr0mk4/test_api/internal/database"
	"github.com/hr0mk4/test_api/internal/models"
	"gorm.io/gorm"
)

func SendCoinHandler(c *gin.Context) {
	dbInstance := database.GetDB()
	var req struct {
		ToUser string `json:"toUser" binding:"required"`
		Amount int    `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Invalid data"})
		return
	}

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "The user is unauthorized"})
		return
	}
	sender, ok := userInterface.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Failed to cast user type"})
		return
	}
	var receiver models.User
	if err := dbInstance.Where("username = ?", req.ToUser).First(&receiver).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "The user does not exist"})
		return
	}
	if sender.ID == receiver.ID {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Can not transfer to yourself"})
		return
	}

	err := dbInstance.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.User{}).Where("id = ? AND coins >= ?", sender.ID, req.Amount).
			Update("coins", gorm.Expr("coins - ?", req.Amount))
		if result.Error != nil || result.RowsAffected == 0 {
			return fmt.Errorf("not enough coins")
		}
		if err := tx.Model(&models.User{}).Where("id = ?", receiver.ID).
			Update("coins", gorm.Expr("coins + ?", req.Amount)).Error; err != nil {
			return err
		}
		transaction := models.Transaction{
			SenderID:   sender.ID,
			ReceiverID: receiver.ID,
			Amount:     req.Amount,
			CreatedAt:  time.Now(),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successful transfer"})
}
