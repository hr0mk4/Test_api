package main

import (
	"context"
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
)

func main() {
	err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	db := database.GetDB()

	r := gin.Default()
	r.POST("/api/auth", handlers.AuthHandler)
	authorized := r.Group("/api")
	authorized.Use(middleware.AuthMiddleware(db))
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
