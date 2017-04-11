package janitor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceContains(t *testing.T) {
	slice := make([]string, 0)
	slice = append(slice, "foobar")
	slice = append(slice, "barfoo")

	assert.True(t, SliceContains(slice, "foobar"))
	assert.False(t, SliceContains(slice, "foobar1"))
}
