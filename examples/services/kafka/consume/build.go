package kafka

import (
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

func Build(
	requestID meta.RequestId,
	rawHeaders service.RawHeaders,
	message service.RawBody,
	rawParams service.RawParams,
) (
	service.MessageTPL, error,
) {

	logger.Printf("Kafka Service: Headers - %+v", rawHeaders)
	logger.Printf("Kafka Service: Params - %+v", rawParams)
	logger.Printf("Kafka Service: Message - %+v", string(message))

	// create service instance
	s := kafkaService{}

	// messages values are in the body
	service.ParseBody(message, &s)

	logger.Printf("Kafka Service: Parsed - %+v", s)

	return &s, nil
}
