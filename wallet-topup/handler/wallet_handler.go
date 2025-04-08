package handler

import (
	"net/http"
	"time"

	"wallet-topup/model"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	svc    model.WalletService
	logger model.Logger
}

func NewWalletHandler(svc model.WalletService, logger model.Logger) *WalletHandler {
	return &WalletHandler{
		svc:    svc,
		logger: logger,
	}
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

	user, err := h.svc.GetUserByID(req.UserID)
	if err != nil || user == nil {
		h.logger.Error("user not found:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	txn, err := h.svc.VerifyTransaction(c.Request.Context(), req.UserID, req.Amount, req.PaymentMethod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction_id": txn.TransactionID,
		"user_id":        txn.UserID,
		"amount":         txn.Amount,
		"payment_method": txn.PaymentMethod,
		"status":         txn.Status,
		"expires_at":     txn.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *WalletHandler) Confirm(c *gin.Context) {
	var req struct {
		TransactionID string `json:"transaction_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	txn, err := h.svc.ConfirmTransaction(c.Request.Context(), req.TransactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.GetUserByID(txn.UserID)
	if err != nil || user == nil {
		h.logger.Error("failed to fetch user balance:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction_id": txn.TransactionID,
		"user_id":        txn.UserID,
		"amount":         txn.Amount,
		"status":         txn.Status,
		"balance":        user.Balance,
	})
}
