package config

import (
	"log"
	"wallet-topup/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase() *gorm.DB {
	//dsn := "host=localhost user=postgres password=pgadmin dbname=wallet-topup port=5432 sslmode=disable"
	dsn := GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	db.AutoMigrate(&model.User{}, &model.Transaction{})
	return db
}
