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
	"gorm.io/gorm"
)

type WalletService struct {
	db     *gorm.DB
	redis  *redis.Client
	ctx    context.Context
	logger *logs.Logger
}

func NewWalletService(db *gorm.DB, redis *redis.Client, logger *logs.Logger) *WalletService {
	return &WalletService{db: db, redis: redis, ctx: context.Background(), logger: logger}
}

func (s *WalletService) VerifyTransaction(userID uint, amount float64, method string) (*model.Transaction, error) {
	var user model.User
	const MaxTopUpAmount = 100000.00

	if err := s.db.First(&user, userID).Error; err != nil {
		s.logger.Error("user not found in DB, user_id=", userID)
		return nil, errors.New("user not found")
	}
	if amount <= 0 {
		s.logger.Warnf("invalid amount %.2f for user_id=%d", amount, userID)
		return nil, errors.New("amount must be greater than zero")
	}
	if amount > MaxTopUpAmount {
		s.logger.Warnf("invalid amount %.2f for user_id=%d", amount, userID)
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

	if err := s.db.Create(txn).Error; err != nil {
		s.logger.Error("failed to create transaction: ", err)
		return nil, err
	}

	data, _ := json.Marshal(txn)
	s.redis.Set(s.ctx, "txn:"+txn.TransactionID, data, 15*time.Minute)

	s.logger.Infof("transaction verified: %s for user %d", txn.TransactionID, userID)
	return txn, nil
}

func (s *WalletService) ConfirmTransaction(transactionID string) (*model.Transaction, error) {
	val, err := s.redis.Get(s.ctx, "txn:"+transactionID).Result()
	var txn model.Transaction
	if err == nil {
		json.Unmarshal([]byte(val), &txn)
	} else {
		if err := s.db.First(&txn, "transaction_id = ?", transactionID).Error; err != nil {
			s.logger.Error("transaction not found: ", transactionID)
			return nil, errors.New("transaction not found")
		}
	}

	if txn.Status != "verified" || time.Now().After(txn.ExpiresAt) {
		s.logger.Warn("transaction expired or already confirmed: ", transactionID)
		return nil, errors.New("transaction expired or already completed")
	}

	s.db.Model(&txn).Update("status", "completed")
	s.db.Model(&model.User{}).Where("user_id = ?", txn.UserID).
		Update("balance", gorm.Expr("balance + ?", txn.Amount))
	s.redis.Del(s.ctx, "txn:"+transactionID)

	s.logger.Infof("transaction confirmed: %s for user %d", txn.TransactionID, txn.UserID)
	return &txn, nil
}

func (s *WalletService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		s.logger.Error("failed to fetch user by ID: ", id)
		return nil, err
	}
	return &user, nil
}
