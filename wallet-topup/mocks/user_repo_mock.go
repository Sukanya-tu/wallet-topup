package mocks

import (
	"wallet-topup/model"

	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) GetUserByID(userID uint) (*model.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserRepoMock) UpdateUserBalance(userID uint, amount float64) error {
	args := m.Called(userID, amount)
	return args.Error(0)
}
