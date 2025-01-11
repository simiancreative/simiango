package temporal

import (
	"context"
	"os"

	"github.com/simiancreative/simiango/logger"
	"go.temporal.io/sdk/testsuite"
)

func MockService() *testsuite.DevServer {
	logger.Mock()

	ctx := context.Background()
	server, err := testsuite.StartDevServer(ctx, testsuite.DevServerOptions{})
	if err != nil {
		logger.Fatalf("Failed to start Temporal server: %v", err)
	}

	logger.Debugf("Mock Temporal service at: %v", server.FrontendHostPort())

	return server
}

func SetHost(server *testsuite.DevServer) {
	os.Setenv("TEMPORAL_HOST", server.FrontendHostPort())
}
