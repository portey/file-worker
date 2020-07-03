package decorator

import (
	"bufio"
	"context"
	"io"
	"time"

	"golang.org/x/sync/errgroup"
)

type Decorated interface {
	PutFile(ctx context.Context, content io.Reader) error
	GetFile(ctx context.Context, cb func(string, time.Time, io.ReadSeeker) error) error
}

type UniqueMatcher interface {
	Exists(string) (bool, error)
	Add(string) error
}

type UniquerFactory interface {
	NewMatcher() UniqueMatcher
}

type UniqueLineDecorator struct {
	decorated      Decorated
	uniquerFactory UniquerFactory
}

func New(decorated Decorated, uniquerFactory UniquerFactory) *UniqueLineDecorator {
	return &UniqueLineDecorator{decorated: decorated, uniquerFactory: uniquerFactory}
}

func (u *UniqueLineDecorator) PutFile(ctx context.Context, content io.Reader) error {
	r, w := io.Pipe()

	group, innerCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		defer w.Close()

		uniquer := u.uniquerFactory.NewMatcher()
		sc := bufio.NewScanner(content)
		for sc.Scan() {
			line := sc.Text()

			exists, err := uniquer.Exists(line)
			if err != nil {
				return err
			}
			if exists {
				continue
			}

			if _, err := w.Write([]byte(line)); err != nil {
				return err
			}
			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}

			if err := uniquer.Add(line); err != nil {
				return err
			}
		}

		return sc.Err()
	})
	group.Go(func() error {
		return u.decorated.PutFile(innerCtx, r)
	})

	return group.Wait()
}

func (u *UniqueLineDecorator) GetFile(ctx context.Context, cb func(string, time.Time, io.ReadSeeker) error) error {
	return u.decorated.GetFile(ctx, cb)
}
