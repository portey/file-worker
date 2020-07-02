package config

type Config struct {
	PrettyLogOutput bool
	LogLevel        string

	HealthCheckPort int
	PrometheusPort  int
	APIPort         int

	StorageBasePath string
}
