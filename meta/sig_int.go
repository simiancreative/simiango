package meta

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/simiancreative/simiango/logger"
)

var cleaners []func()
var done context.Context
var exit context.Context
var initialized bool

func AddCleanup(cleaner func()) {
	cleaners = append(cleaners, cleaner)
}

func CatchSig() (context.Context, context.Context) {
	if initialized {
		return done, exit
	}

	initialized = true

	done, doneCancel := context.WithCancel(context.Background())
	exit, exitCancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigs
		logger.Printf("Meta: begin shutdown (%+v)", sig)

		doneCancel()

		for _, cleaner := range cleaners {
			cleaner()
		}

		logger.Printf("Meta: shutdown complete")

		exitCancel()
	}()

	return done, exit
}
