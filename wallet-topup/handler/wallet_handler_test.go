package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"wallet-topup/handler"
	"wallet-topup/model"
)

func setupHandler() *handler.WalletHandler {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{}, &model.Transaction{})
	db.Create(&model.User{UserID: 1, Balance: 100.0})

	redisClient, _ := redismock.NewClientMock()

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	return handler.NewWalletHandler(db, redisClient, logger)
}

func TestVerifyHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := setupHandler()
	router := gin.New()
	router.POST("/wallet/verify", h.Verify)

	payload := map[string]interface{}{
		"user_id":        1,
		"amount":         50.0,
		"payment_method": "credit_card",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/wallet/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK but got %d", w.Code)
	}
}

func TestConfirmHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{}, &model.Transaction{})
	db.Create(&model.User{UserID: 1, Balance: 100.0})

	txn := &model.Transaction{
		TransactionID: "txn123",
		UserID:        1,
		Amount:        50.0,
		PaymentMethod: "credit_card",
		Status:        "verified",
		ExpiresAt:     time.Now().Add(10 * time.Minute),
	}
	db.Create(txn)

	redisClient, redisMock := redismock.NewClientMock()
	txnJSON, _ := json.Marshal(txn)
	redisMock.ExpectGet("txn:txn123").SetVal(string(txnJSON))
	redisMock.ExpectDel("txn:txn123").SetVal(1)

	logger := logrus.New()
	logger.SetOutput(io.Discard)

	h := handler.NewWalletHandler(db, redisClient, logger)
	router := gin.New()
	router.POST("/wallet/confirm", h.Confirm)

	payload := map[string]interface{}{
		"transaction_id": "txn123",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/wallet/confirm", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK but got %d", w.Code)
	}
}
func TestConfirmHandler_MissingField(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := setupHandler()
	router := gin.New()
	router.POST("/wallet/confirm", h.Confirm)

	body := []byte(`{}`)

	req := httptest.NewRequest("POST", "/wallet/confirm", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request but got %d", w.Code)
	}
}
