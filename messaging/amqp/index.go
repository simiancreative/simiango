package amqp

var session *Session

func Start() {
	NewSession()
}

func Stop() error {
	return session.Close()
}
