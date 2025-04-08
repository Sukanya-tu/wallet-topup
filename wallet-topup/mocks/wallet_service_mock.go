package mocks

import (
	"context"
	"wallet-topup/model"

	"github.com/stretchr/testify/mock"
)

type WalletServiceMock struct {
	mock.Mock
}

func (m *WalletServiceMock) GetUserByID(userID uint) (*model.User, error) {
	args := m.Called(userID)
	if u := args.Get(0); u != nil {
		return u.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *WalletServiceMock) VerifyTransaction(ctx context.Context, userID uint, amount float64, method string) (*model.Transaction, error) {
	args := m.Called(ctx, userID, amount, method)
	if txn := args.Get(0); txn != nil {
		return txn.(*model.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *WalletServiceMock) ConfirmTransaction(ctx context.Context, transactionID string) (*model.Transaction, error) {
	args := m.Called(ctx, transactionID)
	if txn := args.Get(0); txn != nil {
		return txn.(*model.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}
