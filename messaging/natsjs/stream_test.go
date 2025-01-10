package natsjs_test

import (
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/simiancreative/simiango/messaging/natsjs"
	"github.com/simiancreative/simiango/mocks/nats"
)

func TestStream(t *testing.T) {
	natsjs.SetReconnect(false)

	for _, test := range tests {
		// gets the function name to apply to the test run
		funcValue := reflect.ValueOf(test)
		name := runtime.FuncForPC(funcValue.Pointer()).Name()

		// runs test sequentially so the mock server doesnt trip over itself
		t.Run(name, test)
	}
}

var tests = []func(*testing.T){
	testConnect,
	testConnectFailure,
	testInitStreamNoConnection,
	testMessageSetSubject,
	testMessageSetData,
	testMessageSetDataError,
	testMessagePublish,
	testMessagePublishError,
	testPublishNoConnection,
	testPublishNoData,
	testNewConsumer,
	testNewConsumerFailure,
	testNewConsumerConsume,
	testNewConsumerStop,
	testNewConsumerNoConnection,
	testNewConsumerConsumeFailure,
}

func testConnectFailure(t *testing.T) {
	natsjs.Close()

	assert.PanicsWithValue(
		t,
		"NATS_HOST is not set",
		func() {
			os.Unsetenv("NATS_HOST")
			natsjs.Connect()
		},
		"Connect should panic if NATS_HOST is not set",
	)
}

func testConnect(t *testing.T) {
	defer nats.MockServer()()

	assert.NotPanics(
		t,
		func() {
			natsjs.Connect()
		},
		"Connect should not panic if NATS_HOST is set",
	)
}

func testInitStreamNoConnection(t *testing.T) {
	assert.PanicsWithValue(
		t,
		"Connection is not established",
		func() {
			natsjs.Close()
			natsjs.New().InitStream("test")
		},
		"Connect should panic if connection is not set",
	)
}
