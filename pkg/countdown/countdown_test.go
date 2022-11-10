package countdown

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTimeRemaining(t *testing.T) {
	// Arrange
	tNow := time.Now().UTC()
	expected := countdown{t: 0, h: 0, m: 0, s: 0}

	// Act
	result := getTimeRemaining(tNow)

	// Assert
	assert.Equal(t, expected, result, fmt.Sprintf("Incorrect result. Expected %d got %+v", expected, result))
}
