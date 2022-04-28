package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCleaner(t *testing.T) {
	assert.Equal(t, 0, len(cleaners))

	AddCleanup(func() {})

	assert.Equal(t, 1, len(cleaners))
}

func TestCatchSig(t *testing.T) {
	assert.Equal(t, false, initialized)

	done := CatchSig()

	assert.Equal(t, true, initialized)

	sameDone := CatchSig()

	assert.Equal(t, done, sameDone)
}
