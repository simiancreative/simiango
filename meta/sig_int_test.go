package meta

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCleaner(t *testing.T) {
	assert.Equal(t, 0, len(cleaners))

	AddCleanup(func() {})

	assert.Equal(t, 1, len(cleaners))
}

func TestCatchSig(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	assert.Equal(t, false, initialized)

	done, exit := CatchSig()

	assert.Equal(t, true, initialized)

	CatchSig()

	assert.IsType(t, done, ctx)
	assert.IsType(t, exit, ctx)

	cancel()
}
