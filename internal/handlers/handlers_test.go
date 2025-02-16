package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hr0mk4/test_api/internal/handlers"
	"github.com/hr0mk4/test_api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInfoHandlerAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/info", func(c *gin.Context) {
		user := models.User{
			Username: "testuser",
			Coins:    1000,
			Purchases: []models.Purchase{
				{Item: "t-shirt"},
				{Item: "cup"},
				{Item: "t-shirt"},
			},
			ReceivedTransactions: []models.Transaction{
				{SenderID: 2, Amount: 50},
			},
			SentTransactions: []models.Transaction{
				{ReceiverID: 3, Amount: 30},
			},
		}
		c.Set("user", user)
		handlers.InfoHandler(c)
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/info", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "coins")
	assert.Contains(t, resp, "inventory")
	assert.Contains(t, resp, "coinHistory")
}

func TestPurchaseHandlerSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/api/buy/:item", func(c *gin.Context) {
		user := models.User{
			Username:  "buyer",
			Coins:     1000,
			Purchases: []models.Purchase{},
		}
		c.Set("user", user)
		handlers.PurchaseHandler(c)
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/buy/t-shirt", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "message")
}

func TestPurchaseHandlerInsufficientCoins(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/api/buy/:item", func(c *gin.Context) {
		user := models.User{
			Username:  "buyer",
			Coins:     10,
			Purchases: []models.Purchase{},
		}
		c.Set("user", user)
		handlers.PurchaseHandler(c)
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/buy/powerbank", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSendCoinHandlerSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/api/sendCoin", func(c *gin.Context) {
		user := models.User{
			Username: "sender",
			Coins:    500,
		}
		c.Set("user", user)
		handlers.SendCoinHandler(c)
	})

	body := map[string]interface{}{
		"toUser": "receiver",
		"amount": 100,
	}
	jsonBody, _ := json.Marshal(body)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "message")
}

func TestSendCoinHandlerInsufficientFunds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/api/sendCoin", func(c *gin.Context) {
		// Пользователь-отправитель с недостаточным балансом
		user := models.User{
			Username: "sender",
			Coins:    50,
		}
		c.Set("user", user)
		handlers.SendCoinHandler(c)
	})

	body := map[string]interface{}{
		"toUser": "receiver",
		"amount": 100,
	}
	jsonBody, _ := json.Marshal(body)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	// При недостатке средств ожидаем статус 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
