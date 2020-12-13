package service

import (
	"encoding/json"

	"github.com/simiancreative/simiango/meta"
)

type TPL interface {
	Result() (interface{}, error)
}

type Config struct {
	Path   string
	Method string
	Build  func(meta.RequestId, RawBody, RawParams) (TPL, error)
}
type Collection []Config

type RawBody []byte
type RawParams []RawParam
type RawParam struct {
	Key    string
	Value  string
	Values []string
}

func (ps RawParams) Get(name string) (*RawParam, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			return &entry, true
		}
	}
	return nil, false
}

func ParseBody(data []byte, destination interface{}) error {
	return json.Unmarshal(data, destination)
}
