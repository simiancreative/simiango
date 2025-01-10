package nats

import (
	"fmt"
	"net"
	"os"

	natsserver "github.com/nats-io/nats-server/v2/test"
)

func MockServer() func() {
	port, err := getPort()
	if err != nil {
		panic(err)
	}

	os.Setenv("NATS_HOST", fmt.Sprintf("localhost:%v", port))

	opts := natsserver.DefaultTestOptions
	opts.Port = port
	opts.JetStream = true

	return natsserver.RunServer(&opts).Shutdown
}

func getPort() (int, error) {
	addr, _ := net.ResolveTCPAddr("tcp", "localhost:0")

	listener, _ := net.ListenTCP("tcp", addr)
	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port, nil
}
