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
	// implement a read method that streams Datas json representation

	// convert the Data struct to json
	b, err := json.Marshal(e)
	if err != nil {
		return 0, err
	}

	// copy the json to the byte slice
	n = copy(p, b[e.pos:])
	e.pos += n

	// if we have read all the json, return EOF
	if e.pos >= len(b) {
		return n, io.EOF
	}

	return n, nil
}
