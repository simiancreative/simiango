package kafka

import (
	"context"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func BatchMessages(ctx context.Context, values <-chan kafka.Message, maxItems int, maxTimeout time.Duration) chan []kafka.Message {
	batches := make(chan []kafka.Message)

	go func() {
		defer close(batches)

		for keepGoing := true; keepGoing; {
			var batch []kafka.Message
			expire := time.After(maxTimeout)
			for {
				select {
				case value, ok := <-values:
					if !ok {
						keepGoing = false
						goto done
					}

					batch = append(batch, value)
					if len(batch) == maxItems {
						goto done
					}

				case <-expire:
					goto done
				case <-ctx.Done():
					kl.Printf("batch context done")
					return
				}

			}

		done:
			if len(batch) > 0 {
				batches <- batch
			}
		}
	}()

	return batches
}
