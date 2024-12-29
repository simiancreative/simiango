package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/simiancreative/simiango/errors"
	"github.com/simiancreative/simiango/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var exit = os.Exit

func MockExit() {
	exit = func(_ int) {}
}

func New(name string) Root {
	return Root{
		Cmd: &cobra.Command{
			Use: name,
		},
	}
}

type Root struct {
	Cmd *cobra.Command
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func (r Root) Execute() {
	if err := r.Cmd.Execute(); err != nil {
		printStackTrace(err)

		sentry.CaptureException(err)
		sentry.Flush(time.Second * 5)

		exit(1)
	}
}

func Recover() {
	val := recover()
	if val == nil {
		return
	}

	err, ok := val.(error)
	if ok {
		err = errors.Wrap(err, "recovered")
	}

	if !ok {
		err = errors.Wrap(fmt.Errorf("%v", val), "recovered")
	}

	printStackTrace(err)

	sentry.CurrentHub().Recover(err)
	sentry.Flush(time.Second * 5)

	exit(1)
}

func printStackTrace(err error) {
	if logger.Level() < logrus.DebugLevel {
		return
	}

	logger.Errorf(`

=============== LOGGING STACK TRACE =====================

Error: %+v

=============== END LOGGING STACK TRACE =================

		`, err)
}
