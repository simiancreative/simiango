package kafka

import (
	"fmt"

	kafka "github.com/segmentio/kafka-go"
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

func buildService(
	requestID meta.RequestId,
	config service.Config,
	message kafka.Message,
) (service.MessageTPL, error) {
	headers := service.RawHeaders{}
	params := service.RawParams{}

	for _, header := range message.Headers {
		headers = append(headers, service.ParamItem{
			Key:   header.Key,
			Value: string(header.Value),
		})
	}

	params = append(params, service.ParamItem{
		Key:   "topic",
		Value: message.Topic,
	})

	params = append(params, service.ParamItem{
		Key:   "partition",
		Value: fmt.Sprint(message.Partition),
	})

	params = append(params, service.ParamItem{
		Key:   "offset",
		Value: fmt.Sprint(message.Offset),
	})

	params = append(params, service.ParamItem{
		Key:   "high_water_mark",
		Value: fmt.Sprint(message.HighWaterMark),
	})

	params = append(params, service.ParamItem{
		Key:   "time",
		Value: fmt.Sprint(message.Time.Unix()),
	})

	s, err := config.BuildMessages(
		requestID,
		headers,
		message.Value,
		params,
	)

	if err == nil {
		return s, nil
	}

	return nil, err
}
