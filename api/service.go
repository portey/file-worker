package api

import (
	"context"
	"io"
	"time"
)

type Service interface {
	PutFile(ctx context.Context, content io.Reader) error
	GetFile(ctx context.Context, cb func(string, time.Time, io.ReadSeeker) error) error
}
