package main

import (
	"log"
	"os"
	"wallet-topup/config"
	"wallet-topup/handler"
	"wallet-topup/logs"
	"wallet-topup/middleware"
	"wallet-topup/repository"
	"wallet-topup/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	db := config.SetupDatabase()
	redisClient := config.SetupRedis()
	logger := logs.NewLogger()

	userRepo := repository.NewUserRepo(db)
	txnRepo := repository.NewTransactionRepo(db)

	walletService := service.NewWalletService(txnRepo, userRepo, redisClient, logger)
	walletHandler := handler.NewWalletHandler(walletService, logger)

	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		token := config.GenerateToken()
		c.JSON(200, gin.H{"token": token})
	})

	auth := middleware.JWTAuthMiddleware()

	api := r.Group("/api", auth)
	{
		api.POST("/verify", walletHandler.Verify)
		api.POST("/confirm", walletHandler.Confirm)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
