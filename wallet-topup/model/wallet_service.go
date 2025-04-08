package model

import (
	"context"
)

type WalletService interface {
	GetUserByID(userID uint) (*User, error)
	VerifyTransaction(ctx context.Context, userID uint, amount float64, method string) (*Transaction, error)
	ConfirmTransaction(ctx context.Context, transactionID string) (*Transaction, error)
}

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}
