package consumer

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/spf13/cobra"

	"github.com/simiancreative/simiango/errors"
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/messaging/natsjscm"
	"github.com/simiancreative/simiango/messaging/natsjscon"
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
		Subject:      "test",
	}

	connectionConfig := natsjscm.ConnectionConfig{
		Logger: logger,
		URL:    "nats://localhost:4222",
	}
	connector, err := natsjscm.NewConnectionManager(connectionConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create connection manager")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = natsjscon.
		NewConsumer(config).
		SetLogger(logger).
		SetConnector(connector).
		SetStrategy(&Strategy{}).
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
	return map[jetstream.Msg]natsjscon.ProcessStatus{}
}

type Strategy struct {
}

func (s *Strategy) Setup(ctx context.Context) error {
	return nil
}

func (s *Strategy) Consume(ctx context.Context, workerID int) ([]jetstream.Msg, error) {
	return nil, nil
}
