package kafka

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func varAsInt(name string) (bool, int, error) {
	if value, present := os.LookupEnv(name); present && value != "" {
		floatVar, _, err := big.ParseFloat(value, 10, 0, big.ToNearestEven)
		if err != nil {
			kl.Error(fmt.Sprintf("Error converting %s to int", name), fields{"error": err.Error()})
			return true, 0, err
		}
		intVar, _ := floatVar.Uint64()
		return true, int(intVar), nil
	}
	return false, 0, nil
}

func getKafkaReader(kafkaURL string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")

	topic := os.Getenv("KAFKA_READER_TOPIC")
	groupID := os.Getenv("KAFKA_READER_GROUP")

	config := kafka.ReaderConfig{
		Logger:      kl,
		StartOffset: kafka.LastOffset,
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		MinBytes:    1,    // 1B
		MaxBytes:    10e6, // 10MB
	}

	if present, minBytes, err := varAsInt("READ_MIN_BYTES"); present {
		if err != nil {
			return nil
		}
		config.MinBytes = minBytes
	}

	if present, maxBytes, err := varAsInt("READ_MAX_BYTES"); present {
		if err != nil {
			return nil
		}
		config.MaxBytes = maxBytes
	}

	if present, maxWait, err := varAsInt("READ_MAX_WAIT_MILLISECONDS"); present {
		if err != nil {
			return nil
		}
		config.MaxWait = time.Duration(maxWait) * time.Millisecond
	}

	if present, commitInterval, err := varAsInt("READ_COMMIT_INTERVAL_MILLISECONDS"); present {
		if err != nil {
			return nil
		}
		config.CommitInterval = time.Duration(commitInterval) * time.Millisecond
	}

	return kafka.NewReader(config)
}

func closeReader(reader *kafka.Reader) {
	kl.Printf("closing reader")

	if err := reader.Close(); err != nil {
		kl.Fatal("failed to close reader:", fields{"err": err})
	}
}

func NewConsumer(done, sendCtx context.Context, kafkaURL string, wg *sync.WaitGroup) <-chan kafka.Message {
	reader := getKafkaReader(kafkaURL)

	messages := make(chan kafka.Message)

	go func() {
		defer func() {
			kl.Printf("closing consumer...")
			closeReader(reader)
			close(messages)
			kl.Printf("consumer closed, signaling work group...")
			wg.Done()
		}()

		for {
			m, err := readMessages(done, reader)
			if err != nil {
				return
			}
			select {
			case <-done.Done():
				kl.Printf("done in consumer")
				return
			case <-sendCtx.Done():
				kl.Printf("consumer must stop sending")
				return
			case messages <- m:
			}
		}
	}()
	return messages
}

func readMessages(ctx context.Context, reader *kafka.Reader) (kafka.Message, error) {
	m, err := reader.ReadMessage(ctx)
	if err != nil {
		kl.Error("Error reading message", fields{"err": err.Error()})
		return kafka.Message{}, err
	}

	kl.Printf("read message (%+v)", MessageSimplified(m))
	return m, nil
}
