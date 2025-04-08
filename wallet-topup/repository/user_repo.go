package repository

import (
	"wallet-topup/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) UpdateUserBalance(userID uint, amount float64) error {
	return r.DB.Model(&model.User{}).Where("user_id = ?", userID).Update("balance", gorm.Expr("balance + ?", amount)).Error
}
