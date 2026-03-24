package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost		string `mapstructure:"DB_HOST"`
	DBPort		string `mapstructure:"DB_PORT"`
	DBUser		string `mapstructure:"DB_USER"`
	DBPassword	string `mapstructure:"DB_PASSWORD"`
	DBName		string `mapstructure:"DB_NAME"`
	ServerPort	int    `mapstructure:"SERVER_PORT"`
	JWTSecret         string `mapstructure:"JWT_SECRET"`
	JWTExpirationHours int   `mapstructure:"JWT_EXPIRATION_HOURS"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read .env file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if config.DBHost == "" || config.DBPort == "" || config.DBUser == "" || config.DBPassword == "" || config.DBName == "" {
		return nil, errors.New("missing required database configuration fields in environment variables")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("missing JWT_SECRET in environment variables")
	}

	if config.ServerPort == 0 {
		config.ServerPort = 8080
	}

	return &config, nil
}