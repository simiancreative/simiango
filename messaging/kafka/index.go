package kafka

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Start() <-chan bool {
	url := os.Getenv("KAFKA_BROKERS")
	readerTopic := os.Getenv("KAFKA_READER_TOPIC")
	groupID := os.Getenv("KAFKA_READER_GROUP")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	done := make(chan bool)
	go func() {
		sig := <-sigs
		fmt.Println("Start, SIGINT received", sig, "closing...")
		done <- true
		close(done)
	}()

	go func() {
		c1 := NewConsumer(url, readerTopic, groupID, done)

		c2, handlerDone := Handle(c1, done)

		if writerTopic, present := os.LookupEnv("KAFKA_WRITER_TOPIC"); present {
			done := NewProducer(url, writerTopic, c2, done)
			fmt.Println(<-done)
		}
		fmt.Println(<-handlerDone)
	}()
	return done
}
