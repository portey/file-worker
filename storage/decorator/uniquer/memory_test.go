package uniquer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryMatcher_Add(t *testing.T) {
	matcher := New()
	err := matcher.Add("test")
	assert.Nil(t, err)
	err = matcher.Add("tes2")
	assert.Nil(t, err)
	err = matcher.Add("tes2")
	assert.Nil(t, err)

	res, err := matcher.Exists("test")
	assert.Nil(t, err)
	assert.True(t, res)

	res, err = matcher.Exists("tes2")
	assert.Nil(t, err)
	assert.True(t, res)

	res, err = matcher.Exists("test3")
	assert.Nil(t, err)
	assert.False(t, res)
}
