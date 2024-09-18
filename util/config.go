package util

import (
	"errors"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
type Config struct {
	Database DatabaseConfig
	Mongo    MongoConfig
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

// MongoConfig struct for mongodb configuration
type MongoConfig struct {
	Url        string
	Database   string
	Collection string
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
	// Define the flag
	configFileFlag := pflag.StringP("config", "f", "", "Path to the configuration file")
	// Parse the flags
	pflag.Parse()

	// Check if the flag was provided
	if *configFileFlag == "" {
		err = errors.New("flag empty")
		return
	}

	viper.SetConfigFile(*configFileFlag)
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
