package sample

// import "errors"
import (
	"github.com/simiancreative/simiango/data/redis"
	"github.com/simiancreative/simiango/service"
	"github.com/simiancreative/simiango/token"
)

func (s sampleService) Result() (interface{}, error) {
	claims := token.Claims{
		"combat":       80,
		"durability":   45,
		"intelligence": 90,
		"name":         "spider-man",
		"power":        50,
		"speed":        25,
		"strength":     25,
		"tier":         2,
	}

	err := redis.Set("token", token.Gen(claims, 1), 0)
	if err != nil {
		return nil, err
	}

	err = redis.Set("request_id", string(s.id), 0)
	if err != nil {
		return nil, err
	}

	resp := &sampleResp{}

	getErr := redis.Get("request_id", resp)
	if getErr != nil {
		return nil, getErr
	}

	tokenErr := redis.Get("token", resp)
	if tokenErr != nil {
		return nil, tokenErr
	}

	res := service.ToContentResponse([]interface{}{
		claims,
		resp,
		map[string]interface{}{"body": s.body},
		map[string]interface{}{"url_params": s.params},
	}, service.ContentResponseMeta{})

	return res, nil
	// example error response
	// return nil, errors.New("wooble")
}
