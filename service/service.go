package service

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	servertiming "github.com/p768lwy3/gin-server-timing"

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
	Auth(meta.RequestId, RawHeaders, RawBody, RawParams) error
	TPL
}

type ContextTPL interface {
	Result(Req) (interface{}, error)
}

type Message struct {
	Key   string
	Value interface{}
}

type Messages []Message

type MessageTPL interface {
	Result() (Messages, error)
}

type Req struct {
	ID       meta.RequestId
	Headers  RawHeaders
	Body     RawBody
	Params   RawParams
	Receiver interface{}
	Timer    *servertiming.Header
	Context  *gin.Context
}

type Kind int

const (
	DEFAULT Kind = iota // 0
	DIRECT              // 1, and so on.
)

type Config struct {
	Kind            Kind
	IsStream        bool
	IsPrivate       bool
	Key             string
	Path            string
	Method          string
	RequestReceiver func() interface{}
	Build           func(meta.RequestId, RawHeaders, RawBody, RawParams) (TPL, error)
	BuildMessages   func(meta.RequestId, RawHeaders, RawBody, RawParams) (MessageTPL, error)
	Direct          func(req Req) (interface{}, error)
	Auth            func(Req) error
	Before          []func(Config, Req) error
	After           func(Config, Req)
}

type Collection []Config

type RawBody []byte

func ParseBody(data []byte, destination interface{}) error {
	return json.Unmarshal(data, destination)
}
