package main

import (
	"wallet-topup/config"
	"wallet-topup/handler"
	"wallet-topup/logs"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	log := logs.NewLogger()
	db := config.SetupDatabase()
	redis := config.SetupRedis()

	r := gin.Default()
	h := handler.NewWalletHandler(db, redis, log)

	r.POST("/wallet/verify", h.Verify)
	r.POST("/wallet/confirm", h.Confirm)

	log.Info("Server started on :8080")
	r.Run(":8080")
}
