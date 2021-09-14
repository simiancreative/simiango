package stream

import (
	"io"

	"github.com/go-resty/resty/v2"
	"github.com/simiancreative/simiango/service"
)

// https://source.unsplash.com/random/30000x20000

type streamService struct{}

func (s *streamService) Result() (interface{}, error) {
	client := resty.New()

	resp, _ := client.R().
		EnableTrace().
		Get("https://image.tmdb.org/t/p/w780//5nCh9wrr93WdYhquBFEMjZrDgFq.jpg")

	kind := resp.Header()["Content-Type"][0]
	length := resp.Header()["Content-Length"][0]

	writer := func(w io.Writer) bool {
		w.Write(resp.Body())
		return false
	}

	return service.StreamResult{
		Type:   kind,
		Length: length,
		Writer: writer,
	}, nil
}
