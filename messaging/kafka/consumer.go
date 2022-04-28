package kafka

import (
	"context"
	"fmt"
	"strings"

	kafka "github.com/segmentio/kafka-go"
	"github.com/simiancreative/simiango/logger"
)

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		StartOffset: kafka.LastOffset,
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
	})
}

func NewConsumer(kafkaURL, topic, groupID string, done <-chan bool) <-chan kafka.Message {
	/*
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT)
	*/

	reader := getKafkaReader(kafkaURL, topic, groupID)

	out := make(chan kafka.Message)

	/*
		go func() {
			sig := <-sigs
			fmt.Println("Consumer, SIGINT received", sig, "closing...")
			close(out)
			reader.Close()
		}()
	*/

	go func() {
		defer close(out)
		defer reader.Close()
		for {
			select {
			default:
				m, err := reader.ReadMessage(context.Background())
				if err != nil {
					logger.Error("Error reading message", logger.Fields{"err": err.Error()})
					//break
				}
				out <- m
			case <-done:
				fmt.Println("Done in NewConsumer")
				return
			}
		}

	}()

	return out
}
