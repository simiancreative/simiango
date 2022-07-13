package kafka

import (
	"context"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

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

	if minBytes, present := os.LookupEnv("READ_MIN_BYTES"); present {
		floatVar, _, err := big.ParseFloat(minBytes, 10, 0, big.ToNearestEven)
		if err != nil {
			kl.Error("Error converting minBytes to int", fields{"error": err.Error()})
			return nil
		}
		intVar, _ := floatVar.Uint64()
		config.MinBytes = int(intVar)
	}
	if maxBytes, present := os.LookupEnv("READ_MAX_BYTES"); present {
		floatVar, _, err := big.ParseFloat(maxBytes, 10, 0, big.ToNearestEven)
		if err != nil {
			kl.Error("Error converting maxBytes to int", fields{"error": err.Error()})
			return nil
		}
		intVar, _ := floatVar.Uint64()
		config.MaxBytes = int(intVar)
	}
	if maxWait, present := os.LookupEnv("READ_MAX_WAIT_MILLISECONDS"); present {
		floatVar, _, err := big.ParseFloat(maxWait, 10, 0, big.ToNearestEven)
		if err != nil {
			kl.Error("Error converting maxWait to int", fields{"error": err.Error()})
			return nil
		}
		intVar, _ := floatVar.Uint64()
		config.MaxWait = time.Duration(intVar) * time.Millisecond
	}
	if commitInterval, present := os.LookupEnv("READ_COMMIT_INTERVAL"); present {
		floatVar, _, err := big.ParseFloat(commitInterval, 10, 0, big.ToNearestEven)
		if err != nil {
			kl.Error("Error converting maxWait to int", fields{"error": err.Error()})
			return nil
		}
		intVar, _ := floatVar.Uint64()
		config.CommitInterval = time.Duration(intVar) * time.Millisecond
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
