package model

import "time"

type Transaction struct {
	TransactionID string `gorm:"primaryKey;type:uuid"`
	UserID        uint
	Amount        float64 `gorm:"type:numeric(12,2)"`
	PaymentMethod string
	Status        string
	ExpiresAt     time.Time
}

type TransactionRepository interface {
	CreateTransaction(txn *Transaction) error
	GetTransactionByID(transactionID string) (*Transaction, error)
	UpdateTransactionStatus(transactionID string, status string) error
}
