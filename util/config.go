package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all the configurations of the application/
// Values are read by viper from a config file or env variables
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	GrpcServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	ENV                  string        `mapstructure:"ENV"`
	RedisAddress         string        `mapstructure:"REDIS_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Overwrites values if passed in from ENV
	// eg: SERVER_ADDRESS=localhost:8081 make server
	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
