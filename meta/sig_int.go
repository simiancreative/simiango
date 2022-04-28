package meta

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/simiancreative/simiango/logger"
)

var cleaners []func()
var done = make(chan bool, 1)
var initialized bool

func AddCleanup(cleaner func()) {
	cleaners = append(cleaners, cleaner)
}

func CatchSig() chan bool {
	if initialized {
		return done
	}

	initialized = true

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigs
		logger.Printf("Meta: begin shutdown (%+v)", sig)

		close(done)

		for _, cleaner := range cleaners {
			cleaner()
		}

		logger.Printf("Meta: shutdown complete")

		os.Exit(0)
	}()

	return done
}
