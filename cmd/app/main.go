package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hr0mk4/test_api/internal/database"
	"github.com/hr0mk4/test_api/internal/handlers"
	"github.com/hr0mk4/test_api/internal/middleware"
	"github.com/hr0mk4/test_api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=db user=postgres password=secret dbname=merch_store port=5432 sslmode=disable TimeZone=UTC"
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if DB == nil {
		log.Fatalf("no db connection")
	}
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}

	if err := DB.AutoMigrate(&models.User{}, &models.Purchase{}, &models.Transaction{}); err != nil {
		log.Printf("Failed to migrate: %v", err)
		return
	}

	fmt.Println("Database connected")
	database.SetDB(DB)

	r := gin.Default()
	r.POST("/api/auth", handlers.AuthHandler)
	authorized := r.Group("/api")
	authorized.Use(middleware.AuthMiddleware(DB))
	{
		authorized.GET("/info", handlers.InfoHandler)
		authorized.GET("/buy/:item", handlers.PurchaseHandler)
		authorized.POST("/sendCoin", handlers.SendCoinHandler)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		log.Println("Server is running on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server startup error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("The sever has been shut down")
}
