package api

import (
	"io"
	"mime"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const multipartFormDataContentType = "multipart/form-data"

// PUT /v1/file
func (s *Server) putFile(w http.ResponseWriter, r *http.Request) {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-type"))
	if err != nil || contentType != multipartFormDataContentType {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := s.service.PutFile(r.Context(), file); err != nil {
		log.WithError(err).Error("put file err")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GET /v1/file
func (s *Server) getFile(w http.ResponseWriter, r *http.Request) {
	if err := s.service.GetFile(r.Context(), func(name string, lastModified time.Time, content io.ReadSeeker) error {
		http.ServeContent(w, r, name, lastModified, content)
		return nil
	}); err != nil {
		log.WithError(err).Error("get file err")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
