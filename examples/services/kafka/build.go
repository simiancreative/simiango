package kafka

import (
	"fmt"

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

	fmt.Println("Headers", rawHeaders)
	fmt.Println("Params", rawParams)

	// create service instance
	s := kafkaService{}

	// messages values are in the body
	service.ParseBody(message, &s)

	return &s, nil
}
