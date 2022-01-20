package kafka

import (
	"fmt"
	"os"
)

func Read() <-chan []byte {
	url := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_READER_TOPIC")
	groupID := os.Getenv("KAFKA_READER_GROUP")

	return NewConsumer(url, topic, groupID)
}

func Write(in <-chan Result) <-chan bool {
	url := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_WRITER_TOPIC")

	return NewProducer(url, topic, in)
}

func Process(processor Processor) {
	c1 := Read()
	c2 := processEvent(c1, processor)

	done := Write(c2)
	fmt.Println(<-done)
}

func processEvent(c <-chan []byte, processor Processor) <-chan Result {
	out := make(chan Result)

	handler := func(event []byte) {
		result, err := processor(event)
		if err != nil {
			// log error message here
			return
		}

		if result == nil {
			// log error message here
			return
		}

		out <- *result
	}

	go func() {
		for bytes := range c {
			go handler(bytes)
		}
		close(out)
	}()

	return out
}
