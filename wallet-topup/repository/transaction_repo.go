package repository

import (
	"wallet-topup/model"

	"gorm.io/gorm"
)

type TransactionRepo struct {
	DB *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{DB: db}
}

func (r *TransactionRepo) CreateTransaction(txn *model.Transaction) error {
	return r.DB.Create(txn).Error
}

func (r *TransactionRepo) GetTransactionByID(transactionID string) (*model.Transaction, error) {
	var txn model.Transaction
	if err := r.DB.First(&txn, "transaction_id = ?", transactionID).Error; err != nil {
		return nil, err
	}
	return &txn, nil
}

func (r *TransactionRepo) UpdateTransactionStatus(transactionID string, status string) error {
	return r.DB.Model(&model.Transaction{}).Where("transaction_id = ?", transactionID).Update("status", status).Error
}
