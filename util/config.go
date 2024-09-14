package util

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
type Config struct {
	Database DatabaseConfig
	GRPC     GRPCConfig
	Log      Log
}

// DatabaseConfig struct for database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// GRPCConfig struct for GRPC server configuration
type GRPCConfig struct {
	Address string
}

// GRPCConfig struct for GRPC server configuration
type Log struct {
	Level string
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./etc")
	viper.SetConfigName("rd-portfolio")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
