package amqp

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

func buildService(
	requestID meta.RequestId,
	config service.Config,
	d amqp.Delivery,
) (service.TPL, error) {
	parsedHeaders := parseHeaders(d.Headers)
	parsedBody := d.Body
	parsedParams := parseParams(d)

	s, err := config.Build(requestID, parsedHeaders, parsedBody, parsedParams)
	if err == nil {
		return s, nil
	}

	return nil, err
}

func parseHeaders(headers amqp.Table) service.RawHeaders {
	var parsedHeaders = service.RawHeaders{}

	for key, value := range headers {
		parsedHeaders = append(parsedHeaders, service.RawHeader{
			Key:    key,
			Values: []string{fmt.Sprintf("%v", value)},
		})
	}

	return parsedHeaders
}

func parseParams(d amqp.Delivery) service.RawParams {
	var parsedParams = service.RawParams{}

	parsedParams = append(parsedParams, service.RawParam{
		Key:   "ContentType",
		Value: fmt.Sprintf("%v", d.ContentType),
	})

	parsedParams = append(parsedParams, service.RawParam{
		Key:   "ContentEncoding",
		Value: fmt.Sprintf("%v", d.ContentEncoding),
	})

	parsedParams = append(parsedParams, service.RawParam{
		Key:   "DeliveryMode",
		Value: fmt.Sprintf("%v", d.DeliveryMode),
	})

	parsedParams = append(parsedParams, service.RawParam{
		Key:   "Priority",
		Value: fmt.Sprintf("%v", d.Priority),
	})

	parsedParams = append(parsedParams, service.RawParam{
		Key:   "CorrelationId",
		Value: fmt.Sprintf("%v", d.CorrelationId),
	})

	parsedParams = append(parsedParams, service.RawParam{
		Key:   "ReplyTo",
		Value: fmt.Sprintf("%v", d.ReplyTo),
	})

	parsedParams = append(parsedParams, service.RawParam{
		Key:   "Expiration",
		Value: fmt.Sprintf("%v", d.Expiration),
	})

	parsedParams = append(parsedParams, service.RawParam{
		Key:   "MessageId",
		Value: fmt.Sprintf("%v", d.MessageId),
	})

	parsedParams = append(parsedParams, service.RawParam{
		Key:   "Timestamp",
		Value: fmt.Sprintf("%v", d.Timestamp),
	})

	parsedParams = append(parsedParams, service.RawParam{
		Key:   "Type",
		Value: fmt.Sprintf("%v", d.Type),
	})

	parsedParams = append(parsedParams, service.RawParam{
		Key:   "UserId",
		Value: fmt.Sprintf("%v", d.UserId),
	})

	parsedParams = append(parsedParams, service.RawParam{
		Key:   "AppId",
		Value: fmt.Sprintf("%v", d.AppId),
	})

	return parsedParams
}
