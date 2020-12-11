package service

import (
	"encoding/json"

	"simian/context"
)

type TPL interface {
	Result() (interface{}, error)
}

type Config struct {
	Path   string
	Method string
	Build  func(context.RequestId, RawBody, RawParams) (TPL, error)
}
type Collection []Config

type RawBody []byte
type RawParams []RawParam
type RawParam struct {
	Key   string
	Value interface{}
}

func (ps RawParams) Get(name string) (interface{}, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			return entry.Value, true
		}
	}
	return "", false
}

func Unmarshal(data []byte, destination interface{}) error {
	return json.Unmarshal(data, destination)
}
