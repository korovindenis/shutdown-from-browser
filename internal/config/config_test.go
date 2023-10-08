package config

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFileNotExists(t *testing.T) {
	// Arrange
	os.Setenv("CONFIG_PATH", "nonexistent.yml")

	// Act
	cfg, err := New()

	// Assert
	expectedErr := errors.New("config file does not exist: nonexistent.yml")
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, cfg)
}
