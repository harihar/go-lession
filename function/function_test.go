package function

import (
	"testing"

	"github.com/stretcher/testify/assert"
)

func TestAddShouldAddTwoNumbers(t *testing.T) {
	assert.Equal(t, 5, Add(2, 3))
}

func TestSubShouldReturnTheDiff(t *testing.T) {
	assert.Equal(t, 2, Sub(3, 1))
}
