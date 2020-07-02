package healthcheck

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type (
	Server struct {
		http      *http.Server
		checks    []Check
		runErr    error
		readiness bool
	}

	Check func() error
)

func CreateAndRun(port int, checks []Check) *Server {
	service := &Server{
		http: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
		checks: checks,
	}

	service.setupHandlers()
	service.run()

	return service
}

func (s *Server) setupHandlers() {
	handler := http.NewServeMux()

	handler.HandleFunc("/health", s.serve)

	s.http.Handler = handler
}

func (s *Server) run() {
	log.Info("health check service: begin run")

	go func() {
		log.Debug("health check service: addr=", s.http.Addr)
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.runErr = err
			log.WithError(err).Error("health check service: end run")
		}
	}()

	s.readiness = true
}

func (s *Server) Close(ctx context.Context) {
	if err := s.http.Shutdown(ctx); err != nil {
		log.WithError(err).Error("health check service shutdown")
	}
	log.Info("health check service: stopped")
}

func (s *Server) serve(w http.ResponseWriter, r *http.Request) {
	errs := make([]error, 0)
	for _, check := range s.checks {
		if err := check(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		w.WriteHeader(http.StatusInternalServerError)

		for _, err := range errs {
			_, errWrite := w.Write([]byte(fmt.Sprintf("%s\n", err.Error())))
			if errWrite != nil {
				log.Errorf("health check response write error: %s", errWrite.Error())
			}
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
