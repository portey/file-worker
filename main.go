package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/portey/file-worker/api"
	"github.com/portey/file-worker/config"
	"github.com/portey/file-worker/healthcheck"
	"github.com/portey/file-worker/prometheus"
	"github.com/portey/file-worker/storage"
	"github.com/portey/file-worker/storage/decorator"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.ReadOS()
	initLogger(cfg.LogLevel, cfg.PrettyLogOutput)

	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)

	if err := run(ctx, cfg); err != nil {
		cancel()
		log.WithError(err).Fatal("closing service with error")
	}
}

func run(ctx context.Context, cfg config.Config) error {
	store := storage.New(cfg.StorageBasePath)
	uniqueDecorator := decorator.New(store, decorator.NewMemoryUniquerFactory())

	apiSrv := api.CreateAndRun(cfg.APIPort, uniqueDecorator)
	defer closeWithTimeout(apiSrv.Close, 5*time.Second)

	prometheusSrv := prometheus.CreateAndRun(cfg.PrometheusPort)
	defer closeWithTimeout(prometheusSrv.Close, 5*time.Second)

	healthCheckSrv := healthcheck.CreateAndRun(cfg.HealthCheckPort, []healthcheck.Check{
		apiSrv.HealthCheck,
		prometheusSrv.HealthCheck,
	})
	defer closeWithTimeout(healthCheckSrv.Close, 5*time.Second)

	<-ctx.Done()

	return nil
}

func closeWithTimeout(close func(context.Context), d time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	close(ctx)
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Println("Got Interrupt signal")
		stop()
	}()
}

func initLogger(logLevel string, pretty bool) {
	if pretty {
		log.SetFormatter(&log.JSONFormatter{})
	}
	log.SetOutput(os.Stderr)

	switch strings.ToLower(logLevel) {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}
