package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type (
	Server struct {
		http      *http.Server
		service   Service
		runErr    error
		readiness bool
	}
)

func CreateAndRun(port int, service Service) *Server {
	server := &Server{
		http: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
		service: service,
	}

	server.setupHandlers()

	server.run()

	return server
}

func (s *Server) setupHandlers() {
	r := mux.NewRouter()

	r.HandleFunc("/v1/file", s.putFile).Methods("PUT")
	r.HandleFunc("/v1/file", s.getFile).Methods("GET")

	s.http.Handler = r
}

func (s *Server) run() {
	log.Info("api service: begin run")

	go func() {
		log.Debug("api service: addr=", s.http.Addr)
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.runErr = err
			log.WithError(err).Error("api service: end run")
		}
	}()

	s.readiness = true
}

func (s *Server) Close(ctx context.Context) {
	if err := s.http.Shutdown(ctx); err != nil {
		log.WithError(err).Error("api service shutdown")
	}
	log.Info("api service: stopped")
}

func (s *Server) HealthCheck() error {
	if !s.readiness {
		return errors.New("api service is't ready yet")
	}
	if s.runErr != nil {
		return errors.New("api service: run issue")
	}
	if s.service == nil {
		return errors.New("api service: service nil")
	}
	return nil
}
