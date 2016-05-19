package xlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXRequestID(t *testing.T) {
	assert := assert.New(t)

	l := NewLogger()
	assert.NotNil(l)

	assert.NotEmpty(l.XRequestID())
}

func TestNewLogger(t *testing.T) {
	assert := assert.New(t)
	l := NewLogger()
	assert.NotNil(l)
}
