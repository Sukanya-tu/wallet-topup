package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"wallet-topup/logs"
	"wallet-topup/model"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type WalletService struct {
	txnRepo  model.TransactionRepository
	userRepo model.UserRepository
	redis    RedisClient
	logger   logs.Logger
}

func NewWalletService(
	txnRepo model.TransactionRepository,
	userRepo model.UserRepository,
	redis RedisClient,
	logger logs.Logger,
) model.WalletService {
	return &WalletService{
		txnRepo:  txnRepo,
		userRepo: userRepo,
		redis:    redis,
		logger:   logger,
	}
}

func (s *WalletService) VerifyTransaction(ctx context.Context, userID uint, amount float64, method string) (*model.Transaction, error) {
	const MaxTopUpAmount = 100000.00

	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		s.logger.Error("user not found:", userID)
		return nil, errors.New("user not found")
	}

	if amount <= 0 {
		s.logger.Warnf("invalid amount %.2f for user_id=%d", amount, userID)
		return nil, errors.New("amount must be greater than zero")
	}
	if amount > MaxTopUpAmount {
		s.logger.Warnf("amount %.2f exceeds limit for user_id=%d", amount, userID)
		return nil, errors.New("amount exceeds maximum allowed")
	}

	txn := &model.Transaction{
		TransactionID: uuid.New().String(),
		UserID:        userID,
		Amount:        amount,
		PaymentMethod: method,
		Status:        "verified",
		ExpiresAt:     time.Now().Add(15 * time.Minute),
	}

	if err := s.txnRepo.CreateTransaction(txn); err != nil {
		s.logger.Error("failed to create transaction:", err)
		return nil, err
	}

	data, _ := json.Marshal(txn)
	if s.redis != nil {
		s.redis.Set(ctx, "txn:"+txn.TransactionID, data, 15*time.Minute)
	}

	s.logger.Infof("transaction verified: %s", txn.TransactionID)
	return txn, nil
}

func (s *WalletService) ConfirmTransaction(ctx context.Context, transactionID string) (*model.Transaction, error) {
	var val string
	var err error
	if s.redis != nil {
		val, err = s.redis.Get(ctx, "txn:"+transactionID).Result()
	}

	var txn model.Transaction
	if err == nil && val != "" {
		json.Unmarshal([]byte(val), &txn)
	} else {
		dbTxn, err := s.txnRepo.GetTransactionByID(transactionID)
		if err != nil {
			s.logger.Error("transaction not found:", transactionID)
			return nil, errors.New("transaction not found")
		}
		txn = *dbTxn
	}

	if txn.Status != "verified" || time.Now().After(txn.ExpiresAt) {
		s.logger.Warn("transaction expired or already confirmed:", transactionID)
		return nil, errors.New("transaction expired or already completed")
	}

	if err := s.txnRepo.UpdateTransactionStatus(transactionID, "completed"); err != nil {
		s.logger.Error("update status error:", err)
		return nil, err
	}
	if err := s.userRepo.UpdateUserBalance(txn.UserID, txn.Amount); err != nil {
		s.logger.Error("update balance error:", err)
		return nil, err
	}

	if s.redis != nil {
		s.redis.Del(ctx, "txn:"+transactionID)
	}
	s.logger.Infof("transaction confirmed: %s", transactionID)
	txn.Status = "completed"
	return &txn, nil
}

func (s *WalletService) GetUserByID(userID uint) (*model.User, error) {
	return s.userRepo.GetUserByID(userID)
}
