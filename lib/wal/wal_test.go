package wal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	log, err := New()
	assert.Nil(t, err)
	assert.NotNil(t, log)
}
