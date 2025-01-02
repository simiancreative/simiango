package cli_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/simiancreative/simiango/cli"
	"github.com/simiancreative/simiango/config"
	"github.com/simiancreative/simiango/mocks/server"
	"github.com/simiancreative/simiango/monitoring/sentry"
	"github.com/stretchr/testify/assert"
)

var id = uuid.New().String()

var tests = []struct {
	name     string
	param    interface{}
	expected string
}{
	{
		name:     "as string",
		param:    id,
		expected: fmt.Sprintf(`"value":"recovered: %s"`, id),
	},
	{
		name:     "as error",
		param:    fmt.Errorf("error: %s", id),
		expected: fmt.Sprintf(`"value":"recovered: error: %s"`, id),
	},
}

func TestRecoverAndThrow(t *testing.T) {
	config.SetupTest()
	cli.MockExit()

	for _, tt := range tests {
		server, handledRequests, closer := server.Mock()
		defer closer()

		// Set the DSN to the mock server
		mockDSN := fmt.Sprintf("http://publicKey@%s/1", server.Listener.Addr())
		os.Setenv("SENTRY_DSN", mockDSN)

		t.Run(tt.name, func(_ *testing.T) {
			// Set up the Sentry client
			sentry.Enable()

			go func() {
				defer cli.Recover()
				panic(tt.param)
			}()

			select {
			case request := <-handledRequests:
				assert.Equal(t, "/api/1/envelope/", request.R.URL.Path)
				assert.Contains(t, string(request.B), tt.expected)
			case <-time.Tick(10 * time.Second):
				t.Error("never called")
			}
		})
	}
}
