package amqp

import (
	"errors"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/simiancreative/simiango/logger"
)

// This exports a Session object that wraps this library. It
// automatically reconnects when the connection fails, and
// blocks all pushes until the connection succeeds. It also
// confirms every outgoing message, so none are lost.
// It doesn't automatically ack each message, but leaves that
// to the parent process, since it is usage-dependent.
//
// Try running this in one terminal, and `rabbitmq-server` in another.
// Stop & restart RabbitMQ to see how the queue reacts.
type Session struct {
	connection      *amqp.Connection
	channel         *amqp.Channel
	queue           amqp.Queue
	done            chan bool
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isReady         bool
}

const (
	// When reconnecting to the server after connection failure
	reconnectDelay = 5 * time.Second

	// When setting up the channel after a channel exception
	reInitDelay = 2 * time.Second

	// When resending messages the server didn't confirm
	resendDelay = 5 * time.Second
)

var (
	errNotConnected  = errors.New("not connected to a server")
	errAlreadyClosed = errors.New("already closed: not connected to the server")
	errShutdown      = errors.New("session is shutting down")
)

// New creates a new consumer state instance, and automatically
// attempts to connect to the server.
func NewSession() *Session {
	session := Session{
		done: make(chan bool),
	}
	session.handleReconnect()
	return &session
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (session *Session) handleReconnect() {
	for {
		session.isReady = false
		logger.Printf("AMQP: Attempting to connect")

		conn, err := session.connect()

		if err != nil {
			logger.Error(
				"AMQP: Failed to connect. Retrying... ",
				logger.Fields{"err": err},
			)

			select {
			case <-session.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		if done := session.handleReInit(conn); done {
			break
		}
	}
}

// connect will create a new AMQP connection
func (session *Session) connect() (*amqp.Connection, error) {
	amqpURI := os.Getenv("AMQP_URI")
	conn, err := amqp.Dial(amqpURI)

	if err != nil {
		return nil, fmt.Errorf("Amqp Dial: %s", err)
	}

	session.changeConnection(conn)
	logger.Printf("AMQP: Connection Success")
	return conn, nil
}

// handleReconnect will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (session *Session) handleReInit(conn *amqp.Connection) bool {
	for {
		session.isReady = false

		err := session.init(conn)

		if err != nil {
			logger.Error(
				"AMQP: Failed to initialize channel. Retrying... ",
				logger.Fields{"err": err},
			)

			select {
			case <-session.done:
				return true
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-session.done:
			return true
		case <-session.notifyConnClose:
			logger.Printf("AMQP: Connection closed. Reconnecting...")
			return false
		case <-session.notifyChanClose:
			logger.Printf("AMQP: Channel closed. Re-running init...")
		}
	}
}

// init will initialize channel & declare queue
func (session *Session) init(conn *amqp.Connection) error {
	exchange := os.Getenv("AMQP_EXHANGE_NAME")
	exchangeType := os.Getenv("AMQP_EXHANGE_TYPE")

	ch, err := conn.Channel()

	if err != nil {
		return err
	}

	if err = ch.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	queueName := os.Getenv("AMQP_QUEUE_NAME")
	queueKey := os.Getenv("AMQP_QUEUE_KEY")

	queue, err := ch.QueueDeclare(
		queueName,
		false, // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)

	if err != nil {
		return err
	}

	if err = ch.QueueBind(
		queue.Name, // name of the queue
		queueKey,   // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	session.queue = queue
	session.changeChannel(ch)
	session.isReady = true

	logger.Printf("AMQP: Setup Complete")

	go session.Stream()

	return nil
}

// changeConnection takes a new connection to the queue,
// and updates the close listener to reflect this.
func (session *Session) changeConnection(connection *amqp.Connection) {
	session.connection = connection
	session.notifyConnClose = make(chan *amqp.Error)
	session.connection.NotifyClose(session.notifyConnClose)
}

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (session *Session) changeChannel(channel *amqp.Channel) {
	session.channel = channel
	session.notifyChanClose = make(chan *amqp.Error)
	session.notifyConfirm = make(chan amqp.Confirmation, 1)
	session.channel.NotifyClose(session.notifyChanClose)
	session.channel.NotifyPublish(session.notifyConfirm)
}

// Stream will continuously put queue items on the channel.
// It is required to call delivery.Ack when it has been
// successfully processed, or delivery.Nack when it fails.
// Ignoring this will cause data to build up on the server.
func (session *Session) Stream() error {
	if !session.isReady {
		return errNotConnected
	}

	consumerDone := make(chan error)
	tag := os.Getenv("AMQP_CONSUMER_TAG")

	logger.Debug(
		"Queue bound to Exchange, starting Consumer",
		logger.Fields{"consumer tag": tag},
	)

	deliveries, err := session.channel.Consume(
		session.queue.Name,
		tag,   // Consumer
		false, // Auto-Ack
		false, // Exclusive
		false, // No-local
		false, // No-Wait
		nil,   // Args
	)

	if err != nil {
		return fmt.Errorf("Queue Consume: %s", err)
	}

	go handle(deliveries, consumerDone)

	return <-consumerDone
}

// Close will cleanly shutdown the channel and connection.
func (session *Session) Close() error {
	if !session.isReady {
		return errAlreadyClosed
	}
	err := session.channel.Close()
	if err != nil {
		return err
	}
	err = session.connection.Close()
	if err != nil {
		return err
	}
	close(session.done)
	session.isReady = false
	return nil
}
