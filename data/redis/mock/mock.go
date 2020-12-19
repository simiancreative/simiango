package redismock

import (
	"github.com/alicebob/miniredis"
	"github.com/elliotchance/redismock/v8"
	r "github.com/go-redis/redis/v8"

	"github.com/simiancreative/simiango/data/redis"
)

var MockedClient *redismock.ClientMock

func init() {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := r.NewClient(&r.Options{
		Addr: mr.Addr(),
	})

	MockedClient = redismock.NewNiceMock(client)
	redis.C = MockedClient
}
