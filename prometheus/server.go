package prometheus

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	http      *http.Server
	runErr    error
	readiness bool
}

func CreateAndRun(port int) *Server {
	srv := &Server{
		http: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: handler(),
		},
	}

	srv.run()

	return srv
}

func (s *Server) run() {
	log.Info("prometheus service: begin run")

	go func() {
		log.Debug("prometheus service addr:", s.http.Addr)
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.runErr = err
			log.WithError(err).Error("prometheus service: end run")
		}
	}()

	s.readiness = true
}

func (s *Server) Close(ctx context.Context) {
	if err := s.http.Shutdown(ctx); err != nil {
		log.WithError(err).Error("prometheus service shutdown")
	}
	log.Info("prometheus service: stopped")
}

func handler() http.Handler {
	handler := http.NewServeMux()
	handler.Handle("/metrics", promhttp.Handler())
	return handler
}

func (s *Server) HealthCheck() error {
	if !s.readiness {
		return errors.New("prometheus service is't ready yet")
	}
	if s.runErr != nil {
		return errors.New("run prometheus service issue")
	}
	return nil
}
