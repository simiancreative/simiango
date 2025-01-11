package temporal

import (
	"os"

	"go.temporal.io/sdk/client"
)

func Connect(hosts ...string) Client {
	host := os.Getenv("TEMPORAL_HOST")

	if len(hosts) > 0 {
		host = hosts[0]
	}

	if host == "" {
		panic("TEMPORAL_HOST not set")
	}

	var err error
	c, err := client.Dial(client.Options{
		HostPort: host,
		Logger:   NewCustomLogger(),
	})
	if err != nil {
		panic(err)
	}

	return Client{c}
}
