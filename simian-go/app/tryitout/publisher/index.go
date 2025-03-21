package publisher

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"

	"github.com/simiancreative/simiango/errors"
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/messaging/natsjscm"
	"github.com/simiancreative/simiango/messaging/natsjspub"
	cli "github.com/simiancreative/simiango/simian-go/app/tryitout"
)

var (
	count       int
	streamName  string
	subjectName string
)

var cmd = &cobra.Command{
	Use:   "publisher",
	Short: "publish messages to a nats stream",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	cmd.Flags().IntVarP(&count, "count", "c", 100, "Number of messages to publish")
	cmd.Flags().StringVarP(&streamName, "stream", "s", "test", "Stream name")
	cmd.Flags().StringVarP(&subjectName, "subject", "j", "test.something.batch", "Subject name")

	cli.Cmd.AddCommand(cmd)
}

func run() error {
	logger := logger.New()

	// Configure connection
	connectionConfig := natsjscm.ConnectionConfig{
		Logger: logger,
		URL:    "nats://localhost:4222",
	}

	// Create connection manager
	connector, err := natsjscm.NewConnectionManager(connectionConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create connection manager")
	}

	publisher, err := natsjspub.NewPublisher(natsjspub.Dependencies{
		Connector: connector,
		Logger:    logger,
	}, natsjspub.Config{
		StreamName: streamName,
		Subject:    subjectName,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create publisher")
	}

	// Publish messages in rapid succession
	startTime := time.Now()

	logger.Infof("Publishing %d messages to %s...", count, subjectName)

	for i := 1; i <= count; i++ {
		data := fmt.Sprintf("Test message %d", i)

		// Create headers for tracing
		hdr := nats.Header{}
		hdr.Add("X-Message-ID", strconv.Itoa(i))
		hdr.Add("X-Timestamp", time.Now().Format(time.RFC3339Nano))
		hdr.Add("X-Batch-Total", strconv.Itoa(count))

		// Publish with message options
		msg := &nats.Msg{
			Subject: subjectName,
			Header:  hdr,
			Data:    []byte(data),
		}

		ctx := context.Background()
		_, err := publisher.Publish(ctx, msg)

		if err != nil {
			logger.Errorf("Failed to publish message %d: %v", i, err)
			continue
		}
	}

	duration := time.Since(startTime)
	rate := float64(count) / duration.Seconds()

	logger.Infof("Published %d messages in %v (%.2f msgs/sec)", count, duration, rate)

	return nil
}

