package kafka

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func getBatchTimeout() time.Duration {
	amount := 2000

	if amountStr, present := os.LookupEnv("KAFKA_BATCH_TIMEOUT"); present {
		amount, _ = strconv.Atoi(amountStr)
	}

	return time.Duration(amount) * time.Millisecond
}

func getBatchSize() int {
	size := 100

	if sizeStr, present := os.LookupEnv("KAFKA_BATCH_SIZE"); present {
		size, _ = strconv.Atoi(sizeStr)
	}

	return size
}

func GetWriterTopic() string {
	if topic, present := os.LookupEnv("KAFKA_WRITER_TOPIC"); present {
		return topic
	}
	return ""
}

func MustProduce() bool {
	return GetWriterTopic() != ""
}

func getKafkaWriter(kafkaURL string) *kafka.Writer {
	if topic := GetWriterTopic(); topic != "" {
		return &kafka.Writer{
			Logger:    kl,
			Addr:      kafka.TCP(kafkaURL),
			Topic:     topic,
			Balancer:  &kafka.LeastBytes{},
			BatchSize: getBatchSize(),
		}
	}
	return nil
}

func closeWriter(writer *kafka.Writer) {
	kl.Printf("closing writer")

	if err := writer.Close(); err != nil {
		kl.Fatal("failed to close writer:", fields{"err": err})
	}
}

func NewProducer(done context.Context, cancelSend context.CancelFunc, kafkaURL string, in <-chan kafka.Message, wg *sync.WaitGroup) {
	writer := getKafkaWriter(kafkaURL)
	if writer == nil {
		wg.Done()
		return
	}

	kl.Printf("Producer setup (topic: %v)", writer.Topic)

	go func() {
		batchCtx, batchDone := context.WithCancel(done)

		defer func() {
			kl.Printf("closing producer...")
			batchDone()  // Indicates batch to stop sending
			cancelSend() // Indicates upstream to stop sending
			closeWriter(writer)
			kl.Printf("producer closed, signaling work group...")
			wg.Done()
		}()

		batches := BatchMessages(batchCtx, in, getBatchSize(), getBatchTimeout())

		for messages := range batches {
			if err := writeMessage(done, writer, messages); err != nil {
				return
			}
		}
	}()
}

func writeMessage(ctx context.Context, writer *kafka.Writer, messages []kafka.Message) error {
	err := writer.WriteMessages(
		ctx,
		messages...,
	)

	if err != nil {
		kl.Error("Producer - Failed to write messages", fields{"err": err})
		return err
	}

	kl.Info("Producer - wrote messages to kafka", fields{
		"topic":    writer.Topic,
		"count":    len(messages),
		"messages": fmt.Sprintf("%+v", MessagesSimplified(messages)),
	})

	return nil
}
