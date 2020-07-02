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

type UniqueLineDecorator struct {
	decorated Decorated
}

func New(decorated Decorated) *UniqueLineDecorator {
	return &UniqueLineDecorator{decorated: decorated}
}

func (u *UniqueLineDecorator) PutFile(ctx context.Context, content io.Reader) error {
	r, w := io.Pipe()

	group, innerCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		defer w.Close()

		st := make(map[string]struct{})
		sc := bufio.NewScanner(content)
		for sc.Scan() {
			line := sc.Text()
			if _, ok := st[line]; ok {
				continue
			}

			if _, err := w.Write([]byte(line)); err != nil {
				return err
			}
			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}

			st[line] = struct{}{}
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
