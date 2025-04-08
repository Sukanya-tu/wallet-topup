package mocks

import (
	"wallet-topup/model"

	"github.com/stretchr/testify/mock"
)

type TransactionRepoMock struct {
	mock.Mock
}

func (m *TransactionRepoMock) CreateTransaction(txn *model.Transaction) error {
	args := m.Called(txn)
	return args.Error(0)
}

func (m *TransactionRepoMock) GetTransactionByID(transactionID string) (*model.Transaction, error) {
	args := m.Called(transactionID)
	return args.Get(0).(*model.Transaction), args.Error(1)
}

func (m *TransactionRepoMock) UpdateTransactionStatus(transactionID string, status string) error {
	args := m.Called(transactionID, status)
	return args.Error(0)
}
