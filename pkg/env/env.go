package env

import (
	"time"

	"github.com/spf13/viper"
)


//Config defines the enviroment variable configuration
type Config struct {
	LogLevel        string
	LogFormat       string
	GracefulTimeout time.Duration
	Address string
}


//LoadApplicationConfig loads all the system environment variables
func LoadApplicationConfig() Config {
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_FORMAT", "text")
	viper.SetDefault("ADDRESS", ":9000")
	viper.SetDefault("GRACEFUL_TIMEOUT", 20*time.Second)

	viper.ReadInConfig()
	viper.AutomaticEnv()

	return Config{
		Address:            viper.GetString("ADDRESS"),
		LogLevel:        	viper.GetString("LOG_LEVEL"),
		LogFormat:       	viper.GetString("LOG_FORMAT"),
		GracefulTimeout: 	viper.GetDuration("GRACEFUL_TIMEOUT"),
	}
}