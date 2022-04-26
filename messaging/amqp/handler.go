package amqp

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	// https://pkg.go.dev/github.com/rabbitmq/amqp091-go#Delivery
	for d := range deliveries {
		// for each delivery
		// - select matching service
		// - build service
		// - run result
		// - ack message

		logger.Printf(
			"Amqp: got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)

		service, err := findService(d)
		if err != nil {
			handleError("find service error", err, d)
			continue
		}

		if err := handleService(d, service); err != nil {
			handleError("service failed", err, d)
			continue
		}

		d.Ack(false)
	}

	logger.Info("handle: deliveries channel closed", logger.Fields{})
	done <- nil
}

func findService(d amqp.Delivery) (service.Config, error) {
	for _, config := range services {
		if isMatch(config.Key, d) {
			return config, nil
		}
	}

	return service.Config{}, fmt.Errorf("No service found")
}

func isMatch(key string, d amqp.Delivery) bool {
	var content map[string]interface{}

	json.Unmarshal(d.Body, &content)

	for contentKey := range content {
		if key == contentKey {
			return true
		}
	}

	return false
}

func handleService(d amqp.Delivery, config service.Config) error {
	requestID := meta.Id()

	if d.Headers == nil {
		d.Headers = amqp.Table{}
	}
	d.Headers["X-Request-ID"] = string(requestID)

	s, buildErr := buildService(requestID, config, d)
	if buildErr != nil {
		return buildErr
	}

	result, err := s.Result()
	if err != nil {
		return err
	}

	logger.Printf(
		"Amqp: %dB delivery success: [%v] %q, {%v}",
		len(d.Body),
		d.DeliveryTag,
		d.Body,
		result,
	)

	meta.RescuePanic(requestID, d)

	return nil
}

func handleError(message string, err error, d amqp.Delivery) {
	logger.Error(
		fmt.Sprintf("Amqp: %v", message),
		logger.Fields{
			"err":     err.Error(),
			"length":  len(d.Body),
			"dTag":    d.DeliveryTag,
			"headers": d.Headers,
			"body":    fmt.Sprintf("%q", d.Body),
		},
	)

	d.Ack(false)
}
