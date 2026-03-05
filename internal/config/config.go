package config

import (
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

	return &config, nil
}