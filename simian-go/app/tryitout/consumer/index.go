package consumer

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/spf13/cobra"

	"github.com/simiancreative/simiango/circuitbreaker"
	"github.com/simiancreative/simiango/errors"
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/messaging/natsjscm"
	"github.com/simiancreative/simiango/messaging/natsjscon"
	"github.com/simiancreative/simiango/messaging/natsjsstrategypull"
	"github.com/simiancreative/simiango/sig"

	cli "github.com/simiancreative/simiango/simian-go/app/tryitout"
	"github.com/simiancreative/simiango/simian-go/app/tryitout/connection"
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

	connector := connection.Shared
	breaker := newBreaker()
	strategy := newStrategy(breaker)
	consumer := newConsumer(logger, connector, strategy)

	done, exit := sig.New().Catch()

	err := consumer.Start(done)
	if err != nil {
		return errors.Wrap(err, "failed to start consumer")
	}

	<-exit.Done()

	return nil
}

func newConsumer(
	logger natsjscon.Logger,
	connector natsjscm.Connector,
	strategy natsjscon.ConsumptionStrategy,
) natsjscon.Consumer {
	return *natsjscon.NewConsumer(natsjscon.ConsumerConfig{
		StreamName:   "test",
		ConsumerName: "test",
		Subject:      "test.something.>",
	}).
		SetLogger(logger).
		SetConnector(connector).
		SetStrategy(strategy).
		SetProcessor(processor)
}

func newBreaker() circuitbreaker.Breaker {
	breaker, err := circuitbreaker.NewDefault()
	if err != nil {
		logger.Errorf("failed to create circuit breaker: %v", err)
		panic(err)
	}

	return breaker
}

func newStrategy(breaker circuitbreaker.Breaker) natsjscon.ConsumptionStrategy {
	strategy, err := natsjsstrategypull.New(natsjsstrategypull.Config{
		ConsumerName: "test-consumer",
		Breaker:      breaker,
	})

	if err != nil {
		logger.Errorf("failed to create pull strategy: %v", err)
		panic(err)
	}

	return strategy
}

func processor(
	ctx context.Context,
	msgs []jetstream.Msg,
) map[jetstream.Msg]natsjscon.ProcessStatus {
	log := logger.New()

	log.Infof("processing messages: %d", len(msgs))

	processed := map[jetstream.Msg]natsjscon.ProcessStatus{}

	for _, msg := range msgs {
		processed[msg] = natsjscon.Success
		log.Info("processing message", logger.Fields{
			"data": string(msg.Data()),
			"hdr":  msg.Headers(),
		})
	}

	return processed
}
