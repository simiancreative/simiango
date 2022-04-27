package kafka

import (
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func BatchMessages(values <-chan kafka.Message, maxItems int, maxTimeout time.Duration) chan []kafka.Message {
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
