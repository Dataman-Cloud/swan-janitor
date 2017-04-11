package janitor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	defaultConfig := DefaultConfig()
	assert.Equal(t, "0.0.0.0:80", defaultConfig.ListenAddr)
	assert.Equal(t, "lvh.me", defaultConfig.Domain)
}
