package handler

import (
	"net/http"
	"wallet-topup/logs"
	"wallet-topup/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type WalletHandler struct {
	svc    *service.WalletService
	logger *logs.Logger
}

type ConfirmResponse struct {
	TransactionID string  `json:"transaction_id"`
	UserID        uint    `json:"user_id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	Balance       float64 `json:"balance"`
}

func NewWalletHandler(db *gorm.DB, redis *redis.Client, log *logs.Logger) *WalletHandler {
	svc := service.NewWalletService(db, redis, log)
	return &WalletHandler{svc: svc, logger: log}
}

func (h *WalletHandler) Verify(c *gin.Context) {
	var req struct {
		UserID        uint    `json:"user_id"`
		Amount        float64 `json:"amount"`
		PaymentMethod string  `json:"payment_method"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	txn, err := h.svc.VerifyTransaction(req.UserID, req.Amount, req.PaymentMethod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, txn)
}

func (h *WalletHandler) Confirm(c *gin.Context) {
	var req struct {
		TransactionID string `json:"transaction_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	txn, err := h.svc.ConfirmTransaction(req.TransactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.GetUserByID(txn.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user balance"})
		return
	}

	res := ConfirmResponse{
		TransactionID: txn.TransactionID,
		UserID:        txn.UserID,
		Amount:        txn.Amount,
		Status:        txn.Status,
		Balance:       user.Balance,
	}

	c.JSON(http.StatusOK, res)

}
