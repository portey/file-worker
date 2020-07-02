package config

import (
	"github.com/spf13/viper"
)

func ReadOS() Config {
	viper.AutomaticEnv()

	viper.SetEnvPrefix("APP")

	viper.SetDefault("PRETTY_LOG_OUTPUT", true)
	viper.SetDefault("LOG_LEVEL", "DEBUG")

	viper.SetDefault("API_PORT", 8080)
	viper.SetDefault("HEALTH_CHECK_PORT", 8888)
	viper.SetDefault("PROMETHEUS_PORT", 9100)

	//viper.SetDefault("STORAGE_BASE_PATH", "/storage") //todo
	viper.SetDefault("STORAGE_BASE_PATH", "/Users/portey/go/src/github.com/portey/file-worker/temp")

	return Config{
		PrettyLogOutput: viper.GetBool("PRETTY_LOG_OUTPUT"),
		LogLevel:        viper.GetString("LOG_LEVEL"),

		APIPort:         viper.GetInt("API_PORT"),
		HealthCheckPort: viper.GetInt("HEALTH_CHECK_PORT"),
		PrometheusPort:  viper.GetInt("PROMETHEUS_PORT"),

		StorageBasePath: viper.GetString("STORAGE_BASE_PATH"),
	}
}
