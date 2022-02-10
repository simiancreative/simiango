package param

import (
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

func Build(
	requestID meta.RequestId,
	rawHeaders service.RawHeaders,
	rawBody service.RawBody,
	rawParams service.RawParams,
) (
	service.TPL, error,
) {
	c := Company{}
	i := Item{}

	if err := rawParams.Assign(&i); err != nil {
		return nil, err
	}

	if err := rawHeaders.Assign(&c); err != nil {
		return nil, err
	}

	s := paramService{c, i}

	return &s, nil
}
