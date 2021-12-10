package service

import (
	"encoding/json"
	"io"

	"github.com/simiancreative/simiango/meta"
)

type StreamResult struct {
	Type   string
	Length string
	Writer func(io.Writer) bool
}

type TPL interface {
	Result() (interface{}, error)
}

type PrivateTPL interface {
	Auth(meta.RequestId, RawHeaders, RawBody, RawParams) bool
	TPL
}

type Config struct {
	IsStream  bool
	IsPrivate bool
	Key       string
	Path      string
	Method    string
	Build     func(meta.RequestId, RawHeaders, RawBody, RawParams) (TPL, error)
}
type Collection []Config

type RawBody []byte

func ParseBody(data []byte, destination interface{}) error {
	return json.Unmarshal(data, destination)
}
