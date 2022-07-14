package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	kafka "github.com/segmentio/kafka-go"
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

func buildKafkaMessages(done, sendCtx context.Context, messages service.Messages, out chan<- kafka.Message) int {
	msgcnt := len(messages)
	if msgcnt == 0 {
		return 0
	}

	for _, message := range messages {
		marshalled, _ := json.Marshal(message.Value)
		select {
		case out <- kafka.Message{
			Key:   []byte(message.Key),
			Value: marshalled,
		}:
		case <-done.Done():
			kl.Printf("done in handler")
			return -1
		case <-sendCtx.Done():
			// Producer said to stop sending
			kl.Printf("Handler must stop sending")
			return -2
		}
	}
	return msgcnt
}

func Handle(done, sendCtx context.Context, cancelSend context.CancelFunc, messages <-chan kafka.Message, wg *sync.WaitGroup) <-chan kafka.Message {
	handlerName := os.Getenv("KAFKA_HANDLER")
	mustProduce := MustProduce()

	readerConfig, err := findService(handlerName)

	if err != nil {
		kl.Panic("error finding reader service", fields{
			"err": err.Error(),
		})
	}

	results := make(chan kafka.Message)

	handler := func(message kafka.Message) error {
		requestID := meta.Id()

		defer meta.RescuePanic(requestID, MessageSimplified(message))

		service, err := buildService(requestID, readerConfig, message)
		if err != nil {
			kl.Error("error building service", fields{"err": err.Error()})
			return err
		}

		if service == nil {
			kl.Info("no service returned", fields{"message": MessageAsString(message)})
			return nil
		}

		messages, err := service.Result()
		if err != nil {
			kl.Error("error on exec result", fields{"err": err.Error()})
			return nil // Not going to finish in case of Result errors
		}
		if !mustProduce {
			return nil
		}
		if cnt := buildKafkaMessages(done, sendCtx, messages, results); cnt < 0 {
			return errors.New("stop must stop")
		}

		return nil
	}

	go func() {
		defer func() {
			kl.Printf("closing handler...")
			cancelSend() // Indicate upstream to stop sending
			close(results)
			kl.Printf("handler closed, signaling work group...")
			wg.Done()
		}()

		for message := range messages {
			if err := handler(message); err != nil {
				return
			}
		}
	}()

	return results
}

func findService(key string) (service.Config, error) {
	for _, config := range services {
		if config.Key == key {
			return config, nil
		}
	}

	return service.Config{}, fmt.Errorf("No service found")
}
