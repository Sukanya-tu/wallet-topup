package service_test

import (
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"wallet-topup/model"
	"wallet-topup/service"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{}, &model.Transaction{})
	db.Create(&model.User{UserID: 1, Balance: 100.0})
	return db
}

func setupLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	return logger
}

func TestVerifyTransaction_Success(t *testing.T) {
	db := setupTestDB()
	redisClient, mock := redismock.NewClientMock()
	logger := setupLogger()

	svc := service.NewWalletService(db, redisClient, logger)

	mock.ExpectSet("txn:*", "*", 900*time.Second).SetVal("OK")

	txn, err := svc.VerifyTransaction(1, 50.0, "credit_card")
	assert.NoError(t, err)
	assert.Equal(t, "verified", txn.Status)
}

func TestVerifyTransaction_InvalidAmount(t *testing.T) {
	db := setupTestDB()
	redisClient, _ := redismock.NewClientMock()
	logger := setupLogger()

	svc := service.NewWalletService(db, redisClient, logger)

	_, err := svc.VerifyTransaction(1, -10.0, "credit_card")
	assert.Error(t, err)
}

func TestConfirmTransaction_Success(t *testing.T) {
	db := setupTestDB()
	redisClient, mock := redismock.NewClientMock()
	logger := setupLogger()

	svc := service.NewWalletService(db, redisClient, logger)

	txn := &model.Transaction{
		TransactionID: "abc123",
		UserID:        1,
		Amount:        50.0,
		PaymentMethod: "credit_card",
		Status:        "verified",
		ExpiresAt:     time.Now().Add(10 * time.Minute),
	}
	db.Create(txn)

	data, _ := json.Marshal(txn)
	mock.ExpectGet("txn:abc123").SetVal(string(data))
	mock.ExpectDel("txn:abc123").SetVal(1)

	confirmedTxn, err := svc.ConfirmTransaction("abc123")
	assert.NoError(t, err)
	assert.Equal(t, "completed", confirmedTxn.Status)
}
