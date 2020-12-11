package sample

// import "errors"
import (
	"github.com/simiancreative/simiango/data/redis"
	"github.com/simiancreative/simiango/service"
)

func (s sampleService) Result() (interface{}, error) {
	err := redis.Client.Set("request_id", string(s.id), 0)
	if err != nil {
		return nil, err
	}

	val, getErr := redis.Client.Get("request_id")
	if getErr != nil {
		return nil, getErr
	}

	res := service.ToContentResponse([]interface{}{
		map[string]string{"request_id": *val},
		s.body,
		s.params,
	}, service.ContentResponseMeta{})

	return res, nil
	// example error response
	// return nil, errors.New("wooble")
}
