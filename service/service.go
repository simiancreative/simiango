package service

import (
	"encoding/json"

	"github.com/simiancreative/simiango/meta"
)

type TPL interface {
	Result() (interface{}, error)
}

type PrivateTPL interface {
	Auth(meta.RequestId, RawHeaders, RawBody, RawParams) bool
	TPL
}

type Config struct {
	IsPrivate bool
	Path      string
	Method    string
	Build     func(meta.RequestId, RawHeaders, RawBody, RawParams) (TPL, error)
}
type Collection []Config

type RawBody []byte

func ParseBody(data []byte, destination interface{}) error {
	return json.Unmarshal(data, destination)
}
