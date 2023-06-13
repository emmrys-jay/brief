package config

import (
	"log"

	"brief/utility"

	"github.com/spf13/viper"
)

type Configuration struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	SecretKey  string `mapstructure:"SECRET_KEY"`
	RedisHost  string `mapstructure:"REDIS_HOST"`
	RedisPort  string `mapstructure:"REDIS_PORT"`
	PGHost     string `mapstructure:"PG_HOST"`
	PGPort     string `mapstructure:"PG_PORT"`
	PGDatabase string `mapstructure:"PG_DATABASE"`
	PGUser     string `mapstructure:"PG_USER"`
	PGPassword string `mapstructure:"PG_PASSWORD"`
}

// Setup initialize configuration
var (
	Config *Configuration
)

func Setup() {
	var configuration *Configuration
	logger := utility.NewLogger()

	viper.SetConfigName("sample")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Overwrite file env's from environment
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	Config = configuration
	logger.Info("configurations loading successfully")
}

// GetConfig helps you to get configuration data
func GetConfig() *Configuration {
	return Config
}
