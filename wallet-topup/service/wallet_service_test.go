package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"wallet-topup/mocks"
	"wallet-topup/model"
	"wallet-topup/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupLogger() *mocks.LoggerMock {
	logger := new(mocks.LoggerMock)
	logger.On("Infof", mock.Anything, mock.Anything)
	logger.On("Warnf", mock.Anything, mock.Anything, mock.Anything)
	logger.On("Error", mock.Anything, mock.Anything)
	logger.On("Warn", mock.Anything, mock.Anything)
	return logger
}

func TestVerifyTransaction_Success(t *testing.T) {
	txnRepo := new(mocks.TransactionRepoMock)
	userRepo := new(mocks.UserRepoMock)
	redisMock := new(mocks.RedisMock)
	logger := setupLogger()

	userRepo.On("GetUserByID", uint(1)).Return(&model.User{UserID: 1}, nil)
	txnRepo.On("CreateTransaction", mock.Anything).Return(nil)
	redisMock.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	s := service.NewWalletService(txnRepo, userRepo, redisMock, logger)
	txn, err := s.VerifyTransaction(context.Background(), 1, 100.0, "credit_card")

	assert.NoError(t, err)
	assert.Equal(t, uint(1), txn.UserID)
	assert.Equal(t, 100.0, txn.Amount)
	assert.Equal(t, "verified", txn.Status)
}

func TestVerifyTransaction_UserNotFound(t *testing.T) {
	txnRepo := new(mocks.TransactionRepoMock)
	userRepo := new(mocks.UserRepoMock)
	redisMock := new(mocks.RedisMock)
	logger := setupLogger()

	userRepo.On("GetUserByID", uint(99)).Return((*model.User)(nil), errors.New("user not found"))

	s := service.NewWalletService(txnRepo, userRepo, redisMock, logger)
	_, err := s.VerifyTransaction(context.Background(), 99, 100.0, "credit_card")

	assert.EqualError(t, err, "user not found")
}

func TestVerifyTransaction_InvalidAmount(t *testing.T) {
	txnRepo := new(mocks.TransactionRepoMock)
	userRepo := new(mocks.UserRepoMock)
	redisMock := new(mocks.RedisMock)
	logger := setupLogger()

	userRepo.On("GetUserByID", uint(1)).Return(&model.User{UserID: 1}, nil)

	s := service.NewWalletService(txnRepo, userRepo, redisMock, logger)

	_, err := s.VerifyTransaction(context.Background(), 1, -5.0, "credit_card")
	assert.EqualError(t, err, "amount must be greater than zero")

	_, err = s.VerifyTransaction(context.Background(), 1, 1000000.0, "credit_card")
	assert.EqualError(t, err, "amount exceeds maximum allowed")
}

func TestConfirmTransaction_Success(t *testing.T) {
	txnRepo := new(mocks.TransactionRepoMock)
	userRepo := new(mocks.UserRepoMock)
	logger := setupLogger()

	transactionID := uuid.New().String()
	txn := &model.Transaction{
		TransactionID: transactionID,
		UserID:        1,
		Amount:        100.0,
		Status:        "verified",
		ExpiresAt:     time.Now().Add(10 * time.Minute),
	}

	txnRepo.On("GetTransactionByID", transactionID).Return(txn, nil)
	txnRepo.On("UpdateTransactionStatus", transactionID, "completed").Return(nil)
	userRepo.On("UpdateUserBalance", txn.UserID, txn.Amount).Return(nil)

	svc := service.NewWalletService(txnRepo, userRepo, nil, logger)
	res, err := svc.ConfirmTransaction(context.Background(), transactionID)

	assert.NoError(t, err)
	assert.Equal(t, "completed", res.Status)
	assert.Equal(t, transactionID, res.TransactionID)
}

func TestConfirmTransaction_Expired(t *testing.T) {

	txnRepo := new(mocks.TransactionRepoMock)
	userRepo := new(mocks.UserRepoMock)
	logger := setupLogger()

	transactionID := uuid.New().String()
	expiredTxn := &model.Transaction{
		TransactionID: transactionID,
		UserID:        1,
		Amount:        100.0,
		Status:        "verified",
		ExpiresAt:     time.Now().Add(-10 * time.Minute),
	}

	txnRepo.On("GetTransactionByID", transactionID).Return(expiredTxn, nil)

	svc := service.NewWalletService(txnRepo, userRepo, nil, logger)
	res, err := svc.ConfirmTransaction(context.Background(), transactionID)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, "transaction expired or already completed", err.Error())
}

func TestConfirmTransaction_InvalidStatus(t *testing.T) {
	txnRepo := new(mocks.TransactionRepoMock)
	userRepo := new(mocks.UserRepoMock)
	logger := setupLogger()

	transactionID := uuid.New().String()

	txn := &model.Transaction{
		TransactionID: transactionID,
		UserID:        1,
		Amount:        200.0,
		Status:        "completed",
		ExpiresAt:     time.Now().Add(10 * time.Minute),
	}

	txnRepo.On("GetTransactionByID", transactionID).Return(txn, nil)

	s := service.NewWalletService(txnRepo, userRepo, nil, logger)
	res, err := s.ConfirmTransaction(context.Background(), transactionID)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, "transaction expired or already completed", err.Error())
}

func TestConfirmTransaction_NotFound(t *testing.T) {
	txnRepo := new(mocks.TransactionRepoMock)
	userRepo := new(mocks.UserRepoMock)
	logger := setupLogger()

	transactionID := uuid.New().String()

	txnRepo.On("GetTransactionByID", transactionID).Return((*model.Transaction)(nil), errors.New("not found"))

	s := service.NewWalletService(txnRepo, userRepo, nil, logger)
	res, err := s.ConfirmTransaction(context.Background(), transactionID)

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.EqualError(t, err, "transaction not found")
}
