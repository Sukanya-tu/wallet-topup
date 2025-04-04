package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config.yaml: %v", err)
	}
	log.Println("✅ Loaded config.yaml")
}

func GetDSN() string {
	return "host=" + viper.GetString("database.host") +
		" user=" + viper.GetString("database.user") +
		" password=" + viper.GetString("database.password") +
		" dbname=" + viper.GetString("database.name") +
		" port=" + viper.GetString("database.port") +
		" sslmode=" + viper.GetString("database.sslmode")
}

func GetRedisAddr() string {
	return viper.GetString("redis.host") + ":" + viper.GetString("redis.port")
}
