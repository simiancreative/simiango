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

	err := redis.Client.Set("token", token.Gen(claims, 1), 0)
	if err != nil {
		return nil, err
	}

	err = redis.Client.Set("request_id", string(s.id), 0)
	if err != nil {
		return nil, err
	}

	reqid, getErr := redis.Client.Get("request_id")
	if getErr != nil {
		return nil, getErr
	}

	tokenid, tokenErr := redis.Client.Get("token")
	if tokenErr != nil {
		return nil, tokenErr
	}

	_ = service.ToContentResponse([]interface{}{
		claims,
		map[string]string{"token": *tokenid},
		map[string]string{"request_id": *reqid},
		map[string]interface{}{"body": s.body},
		map[string]interface{}{"url_params": s.params},
	}, service.ContentResponseMeta{})

	return nil, nil
	// example error response
	// return nil, errors.New("wooble")
}
