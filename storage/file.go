package storage

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"
)

const fileName = "file.txt"

type FileStorage struct {
	basePath string
}

func New(basePath string) *FileStorage {
	return &FileStorage{basePath: basePath}
}

func (f *FileStorage) PutFile(_ context.Context, content io.Reader) error {
	file, err := ioutil.TempFile(f.basePath, "processing")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	if _, err := io.Copy(file, content); err != nil {
		return err
	}

	return os.Rename(file.Name(), f.getFilePath())
}

func (f *FileStorage) GetFile(_ context.Context, cb func(string, time.Time, io.ReadSeeker) error) error {
	fi, err := os.Open(f.getFilePath())
	if err != nil {
		return err
	}
	defer fi.Close()

	statInfo, err := fi.Stat()
	if err != nil {
		return err
	}

	return cb(statInfo.Name(), statInfo.ModTime(), fi)
}

func (f *FileStorage) getFilePath() string {
	return path.Join(f.basePath, fileName)
}
