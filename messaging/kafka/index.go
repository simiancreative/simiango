package kafka

import (
	"context"
	"os"
	"sync"
)

func Start() (close func()) {
	done, doneCancel := context.WithCancel(context.Background())

	for _, service := range services {
		kl.Printf("starting service: %s", service.Key)
	}

	url := os.Getenv("KAFKA_BROKERS")

	var wg sync.WaitGroup
	wg.Add(3)

	sendCtx, cancelSend := context.WithCancel(context.Background())

	messages := NewConsumer(done, sendCtx, url, &wg)
	results := Handle(done, sendCtx, cancelSend, messages, &wg)

	NewProducer(done, cancelSend, url, results, &wg)

	return func() {
		doneCancel()
		// wait for the channels to close.
		wg.Wait()

		kl.Printf("cleanup complete")
	}
}
