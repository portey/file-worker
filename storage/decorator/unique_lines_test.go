package decorator

import (
	"context"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/portey/file-worker/storage/decorator/mock"
	"github.com/stretchr/testify/assert"
)

func Test_DecoratedPut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	decorated := mock.NewMockDecorated(ctrl)
	decorated.EXPECT().
		PutFile(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, content io.Reader) error {
			filtered, err := ioutil.ReadAll(content)
			assert.Nil(t, err)

			assert.Equal(t, `
word1
word2
`, string(filtered))

			return nil
		})

	decorator := New(decorated, NewMemoryUniquerFactory())
	ctx := context.Background()

	content := strings.NewReader(`
word1
word1
word2
word1
word2
word1
`)

	err := decorator.PutFile(ctx, content)
	assert.Nil(t, err)
}

func Test_DecoratedGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var f func(string, time.Time, io.ReadSeeker) error
	decorated := mock.NewMockDecorated(ctrl)
	decorated.EXPECT().
		GetFile(gomock.Any(), gomock.Eq(f)).
		Return(nil)

	decorator := New(decorated, NewMemoryUniquerFactory())
	ctx := context.Background()

	err := decorator.GetFile(ctx, f)
	assert.Nil(t, err)

}
