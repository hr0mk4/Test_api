package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hr0mk4/test_api/internal/database"
	"github.com/hr0mk4/test_api/internal/models"
)

func InfoHandler(c *gin.Context) {
	dbInstance := database.GetDB()
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

	if err := dbInstance.Preload("Purchases").
		Preload("SentTransactions").
		Preload("ReceivedTransactions").
		First(&user).Error; err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Failed getting user's data"})
		return
	}

	inventoryMap := make(map[string]int, len(user.Purchases))
	for _, p := range user.Purchases {
		inventoryMap[p.Item]++
	}

	inventory := make([]gin.H, 0, len(inventoryMap))
	for item, quantity := range inventoryMap {
		inventory = append(inventory, gin.H{"item": item, "quantity": quantity})
	}

	// Формирование истории переводов (полученных и отправленных)
	receivedHistory := make([]gin.H, 0, len(user.ReceivedTransactions))
	for _, t := range user.ReceivedTransactions {
		var sender models.User
		if err := dbInstance.First(&sender, t.SenderID).Error; err == nil {
			receivedHistory = append(receivedHistory, gin.H{"fromUser": sender.Username, "amount": t.Amount})
		}
	}

	sentHistory := make([]gin.H, 0, len(user.SentTransactions))
	for _, t := range user.SentTransactions {
		var receiver models.User
		if err := dbInstance.First(&receiver, t.ReceiverID).Error; err == nil {
			sentHistory = append(sentHistory, gin.H{"toUser": receiver.Username, "amount": t.Amount})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"coins":     user.Coins,
		"inventory": inventory,
		"coinHistory": gin.H{
			"received": receivedHistory,
			"sent":     sentHistory,
		},
	})
}
