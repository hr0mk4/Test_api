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

func PurchaseHandler(c *gin.Context) {
	dbInstance := database.GetDB()
	item := c.Param("item")
	price, exists := models.Merch_price[item]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "No such type of merch"})
		return
	}

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "The user is unauthorized"})
		return
	}
	user, ok := userInterface.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Failed to cast user type"})
		return
	}

	if user.Coins < price {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Not enough coins"})
		return
	}

	err := dbInstance.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.User{}).Where("id = ? AND coins >= ?", user.ID, price).
			Update("coins", gorm.Expr("coins - ?", price))
		if result.Error != nil || result.RowsAffected == 0 {
			return fmt.Errorf("operation failed")
		}
		purchase := models.Purchase{
			UserID:    user.ID,
			Item:      item,
			Price:     price,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&purchase).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successful purchase"})
}
