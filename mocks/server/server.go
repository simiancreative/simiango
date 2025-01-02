package server

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/kr/pretty"
)

type HandledRequest struct {
	W http.ResponseWriter
	R *http.Request
	B []byte
}

func Mock() (*httptest.Server, chan HandledRequest, func()) {
	getChans := make(chan HandledRequest)

	// Create a new instance of the server with the desired handler.
	mockServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			body, _ := io.ReadAll(r.Body)
			pretty.Println(string(body))

			getChans <- HandledRequest{w, r, body}
		}),
	)

	closer := func() {
		defer mockServer.Close()
		close(getChans)
	}

	return mockServer, getChans, closer
}
