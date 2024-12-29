package sig

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/simiancreative/simiango/logger"
)

func New() *Sig {
	return &Sig{}
}

type Sig struct {
	cleaners    []func()
	done        context.Context
	exit        context.Context
	Initialized bool
}

func (s *Sig) AddCleanup(cleaner func()) *Sig {
	s.cleaners = append(s.cleaners, cleaner)

	return s
}

func (s *Sig) Catch() (context.Context, context.Context) {
	if s.Initialized {
		return s.done, s.exit
	}

	s.Initialized = true

	done, doneCancel := context.WithCancel(context.Background())
	exit, exitCancel := context.WithCancel(context.Background())

	s.done = done
	s.exit = exit

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		logger.Enable()

		sig := <-sigs
		logger.Printf("Sig: begin shutdown (%+v)", sig)

		doneCancel()

		for _, cleaner := range s.cleaners {
			cleaner()
		}

		logger.Printf("Sig: shutdown complete")

		exitCancel()
	}()

	return done, exit
}
