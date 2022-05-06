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

	// return nil, nil to skip an event without error.
	return nil, nil
}
