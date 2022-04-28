package kafka

import (
	"context"
	"log"

	kafka "github.com/segmentio/kafka-go"
	"github.com/simiancreative/simiango/logger"
)

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func closeWriter(writer *kafka.Writer) {
	logger.Printf("closing writer\n")
	if err := writer.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func NewProducer(kafkaURL, topic string, in <-chan []kafka.Message, done <-chan bool) <-chan bool {
	writer := getKafkaWriter(kafkaURL, topic)
	//done := make(chan bool)

	/*
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT)
		go func() {
			sig := <-sigs
			fmt.Println("Producer, SIGINT received", sig, "closing...")
			done <- true
			closeWriter(writer)
		}()
	*/

	go func() {
		defer closeWriter(writer)

		for messages := range in {
			select {
			case <-done:
				return
			default:
				err := writer.WriteMessages(
					context.Background(),
					messages...,
				)
				if err != nil {
					log.Fatal("Failed to write messages:", err)
				}
			}
		}

		//done <- true
	}()
	return done
}
