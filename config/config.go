package config

import (
	"time"

	"github.com/rs/zerolog"
)

var (
	// access token
	expHours       = 4
	AccessTokenexp = time.Now().Add(time.Hour * time.Duration(expHours)).Unix()

	// refresh token
	expDays         = 30
	RefreshTokenexp = time.Now().Add(time.Hour * 24 * time.Duration(expDays)).Unix()

	// SecretKey for jwt
	SecretKey = []byte("as;lkdjfeopiriour")
)

type LoggerConfig struct {
	LogLevel    zerolog.Level
	LogFilePath string
	MaxSize     int64
	MaxBackups  int
	MaxAge      int
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Address string
	Port    int
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

func Load() (*Config, error) {
	// TODO: Implement loading from environment variables or config file
	return &Config{
		Server: ServerConfig{
			Address: ":8080",
			Port:    8080,
		},
		Database: DatabaseConfig{
			Driver: "sqlite",
			DSN:    "./app.db",
		},
		Logger: LoggerConfig{
			LogLevel:    zerolog.InfoLevel,
			LogFilePath: "./storage/logs/app.log",
			MaxSize:     50,
			MaxBackups:  3,
			MaxAge:      30,
		},
	}, nil
}
