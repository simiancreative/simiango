package consumer

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"

	"github.com/simiancreative/simiango/errors"
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/messaging/natsjscm"
	"github.com/simiancreative/simiango/messaging/natsjscon"
	"github.com/simiancreative/simiango/messaging/natsjsstrategypull"
	"github.com/simiancreative/simiango/sig"
	cli "github.com/simiancreative/simiango/simian-go/app/tryitout"
)

var cmd = &cobra.Command{
	Use:   "consumer",
	Short: "run a nats consumer",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	cli.Cmd.AddCommand(cmd)
}

func run() error {
	logger := logger.New()

	config := natsjscon.ConsumerConfig{
		StreamName:   "test",
		ConsumerName: "test",
		Subject:      "test.something.>",
	}

	connectionConfig := natsjscm.ConnectionConfig{
		Logger: logger,
		URL:    "nats://localhost:4222",
	}
	connector, err := natsjscm.NewConnectionManager(connectionConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create connection manager")
	}

	strategy, err := natsjsstrategypull.New(natsjsstrategypull.Config{
		ConsumerName: "test-consumer",
	})
	if err != nil {
		return errors.Wrap(err, "failed to create pull strategy")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = natsjscon.
		NewConsumer(config).
		SetLogger(logger).
		SetConnector(connector).
		SetStrategy(strategy).
		SetProcessor(processor).
		Start(ctx)

	if err != nil {
		return errors.Wrap(err, "failed to start consumer")
	}

	_, exit := sig.
		New().
		Catch()

	<-exit.Done()

	return nil
}

func processor(
	ctx context.Context,
	msgs []jetstream.Msg,
) map[jetstream.Msg]natsjscon.ProcessStatus {
	logger.Debug("processing messages", logger.Fields{
		"messages": msgs,
	})

	litter.Dump(msgs[0].Data())

	return map[jetstream.Msg]natsjscon.ProcessStatus{}
}
