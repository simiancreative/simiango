package sig_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/simiancreative/simiango/sig"
	"github.com/stretchr/testify/assert"
)

func TestCatchSig(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := sig.New().AddCleanup(func() {})

	assert.Equal(t, false, s.Initialized)

	done, exit := s.Catch()

	assert.Equal(t, true, s.Initialized)

	s.Catch()

	assert.IsType(t, done, ctx)
	assert.IsType(t, exit, ctx)

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		time.Sleep(1 * time.Second)
		proc.Signal(os.Interrupt)
	}()

	_, ok := <-exit.Done()
	assert.Equal(t, false, ok)

	cancel()
}
