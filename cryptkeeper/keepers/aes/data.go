package aes

import (
	"encoding/json"
	"io"
	// "github.com/sanity-io/litter"
)

type Data struct {
	Hash string `json:"hash"`
	Salt string `json:"salt"`
	pos  int
}

func (e Data) Read(p []byte) (n int, err error) {
	data, err := json.Marshal(e)
	if err != nil {
		return
	}

	if e.pos >= len(data) {
		return 0, io.EOF
	}

	n = copy(p, data[e.pos:])
	e.pos += n

	return n, nil
}
