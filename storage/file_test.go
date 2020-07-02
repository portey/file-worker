package storage

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_PutFile(t *testing.T) {
	currentDir, err := os.Getwd()
	assert.Nil(t, err)

	ctx := context.Background()
	store := New(path.Join(currentDir, "testdata"))

	err = store.PutFile(ctx, strings.NewReader("content"))
	assert.Nil(t, err)

	createdFile := path.Join(currentDir, "testdata", fileName)
	content, err := ioutil.ReadFile(createdFile)
	assert.Nil(t, err)

	assert.Equal(t, "content", string(content))
	assert.Nil(t, os.Remove(createdFile))
}

func Test_GetFile(t *testing.T) {
	currentDir, err := os.Getwd()
	assert.Nil(t, err)

	filePath := path.Join(currentDir, "testdata", fileName)
	f, err := os.Create(filePath)
	assert.Nil(t, err)
	defer os.Remove(filePath)

	_, err = f.Write([]byte("content"))
	assert.Nil(t, err)

	ctx := context.Background()
	store := New(path.Join(currentDir, "testdata"))
	err = store.GetFile(ctx, func(s string, time time.Time, seeker io.ReadSeeker) error {
		assert.Equal(t, fileName, s)

		content, err := ioutil.ReadAll(seeker)
		assert.Nil(t, err)
		assert.Equal(t, "content", string(content))

		return nil
	})
	assert.Nil(t, err)
}
