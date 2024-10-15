package util

import (
	"errors"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type CacheType string

// Config stores all configuration of the application.
type Config struct {
	Database DatabaseConfig
	Mongo    MongoConfig
	GRPC     GRPCConfig
	Cache    CacheConfig
	Log      LogConfig
	Encrypt  EncryptConfig
}

// DatabaseConfig struct for database configuration
type DatabaseConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Database   string
	Encryption string
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

// LogConfig struct for logging configuration
// Level specifies the logging level, e.g. "info", "debug", "error".
type LogConfig struct {
	Level string
}

// CacheConfig struct for cache configuration
type CacheConfig struct {
	Type     CacheType
	Host     string
	Port     int
	Username string
	Password string
}

// EncryptConfig holds the configuration for encryption.
// The Key field specifies the encryption key.
type EncryptConfig struct {
	Key string
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
