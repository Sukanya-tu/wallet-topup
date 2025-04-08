package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"wallet-topup/handler"
	"wallet-topup/mocks"
	"wallet-topup/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(h *handler.WalletHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/wallet/verify", h.Verify)
	r.POST("/wallet/confirm", h.Confirm)
	return r
}

func TestVerify_Success(t *testing.T) {
	txnRepo := new(mocks.TransactionRepoMock)
	userRepo := new(mocks.UserRepoMock)
	logger := new(mocks.LoggerMock)

	userRepo.On("GetUserByID", uint(1)).Return(&model.User{
		UserID: 1,
	}, nil)

	txnID := uuid.New().String()
	txn := &model.Transaction{
		TransactionID: txnID,
		UserID:        1,
		Amount:        100.50,
		PaymentMethod: "credit_card",
		Status:        "verified",
		ExpiresAt:     time.Now().Add(15 * time.Minute),
	}

	txnRepo.On("CreateTransaction", txn).Return(nil)

	svc := &mocks.WalletServiceMock{}
	svc.On("GetUserByID", uint(1)).Return(&model.User{UserID: 1}, nil)
	svc.On("VerifyTransaction", mock.Anything, uint(1), 100.50, "credit_card").Return(txn, nil)

	h := handler.NewWalletHandler(svc, logger)
	router := setupRouter(h)

	body := map[string]interface{}{
		"user_id":        1,
		"amount":         100.50,
		"payment_method": "credit_card",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/wallet/verify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, txnID, res["transaction_id"])
	assert.Equal(t, float64(1), res["user_id"])
	assert.Equal(t, 100.50, res["amount"])
	assert.Equal(t, "credit_card", res["payment_method"])
	assert.Equal(t, "verified", res["status"])
	assert.NotEmpty(t, res["expires_at"])
}

func TestVerify_UserNotFound(t *testing.T) {
	logger := new(mocks.LoggerMock)
	svc := new(mocks.WalletServiceMock)

	svc.On("GetUserByID", uint(99)).Return((*model.User)(nil), errors.New("not found"))
	logger.On("Error", "user not found:", mock.Anything)

	h := handler.NewWalletHandler(svc, logger)
	router := setupRouter(h)

	body := map[string]interface{}{
		"user_id":        99,
		"amount":         100,
		"payment_method": "credit_card",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/wallet/verify", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVerify_InvalidAmount(t *testing.T) {
	logger := new(mocks.LoggerMock)
	svc := new(mocks.WalletServiceMock)

	svc.On("GetUserByID", uint(1)).Return(&model.User{UserID: 1}, nil)
	svc.On("VerifyTransaction", mock.Anything, uint(1), 0.0, "credit_card").
		Return(nil, errors.New("amount must be greater than zero"))

	h := handler.NewWalletHandler(svc, logger)
	router := setupRouter(h)

	body := map[string]interface{}{
		"user_id":        1,
		"amount":         0,
		"payment_method": "credit_card",
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/wallet/verify", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestConfirm_Success(t *testing.T) {
	logger := new(mocks.LoggerMock)

	transactionID := uuid.New().String()
	txn := &model.Transaction{
		TransactionID: transactionID,
		UserID:        1,
		Amount:        200,
		Status:        "completed",
	}

	user := &model.User{UserID: 1, Balance: 500.75}

	svc := &mocks.WalletServiceMock{}
	svc.On("ConfirmTransaction", mock.Anything, transactionID).Return(txn, nil)
	svc.On("GetUserByID", uint(1)).Return(user, nil)

	h := handler.NewWalletHandler(svc, logger)
	router := setupRouter(h)

	body := map[string]interface{}{"transaction_id": transactionID}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/wallet/confirm", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, transactionID, res["transaction_id"])
	assert.Equal(t, float64(1), res["user_id"])
	assert.Equal(t, 200.0, res["amount"])
	assert.Equal(t, "completed", res["status"])
	assert.Equal(t, 500.75, res["balance"])
}

func TestConfirm_Expired(t *testing.T) {
	logger := new(mocks.LoggerMock)
	svc := new(mocks.WalletServiceMock)

	txnID := uuid.New().String()

	svc.On("ConfirmTransaction", mock.Anything, txnID).
		Return(nil, errors.New("transaction has expired"))

	h := handler.NewWalletHandler(svc, logger)
	router := setupRouter(h)

	body := map[string]interface{}{"transaction_id": txnID}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/wallet/confirm", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestConfirm_InvalidStatus(t *testing.T) {
	logger := new(mocks.LoggerMock)
	svc := new(mocks.WalletServiceMock)

	txnID := uuid.New().String()

	svc.On("ConfirmTransaction", mock.Anything, txnID).
		Return(nil, errors.New("transaction is already completed"))

	h := handler.NewWalletHandler(svc, logger)
	router := setupRouter(h)

	body := map[string]interface{}{"transaction_id": txnID}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/wallet/confirm", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
