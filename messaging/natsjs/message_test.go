package natsjs_test

import (
	"encoding/json"
	"testing"

	"github.com/simiancreative/simiango/messaging/natsjs"
	"github.com/simiancreative/simiango/mocks/nats"
	"github.com/stretchr/testify/assert"
)

func testMessageSetSubject(t *testing.T) {
	msg := &natsjs.Message{}
	msg.SetSubject("prefix", "suffix")
	assert.Equal(t, "prefix.suffix", msg.Subject)
}

func testMessageSetData(t *testing.T) {
	msg := &natsjs.Message{}
	data := map[string]string{"key": "value"}
	msg.SetData(data)

	expectedData, _ := json.Marshal(data)
	assert.Equal(t, expectedData, msg.Data)
}

func testMessageSetDataError(t *testing.T) {
	assert.Panics(
		t,
		func() {
			msg := &natsjs.Message{}
			data := make(chan int)
			msg.SetData(data)
		},
		"SetData should panic if data is not serializable",
	)
}

func testMessagePublish(t *testing.T) {
	defer nats.MockServer()()
	natsjs.Connect()
	natsjs.SetTimeout(2)

	err := natsjs.New().
		InitStream("test").
		NewMessage().
		SetSubject("test", "subject").
		SetData("test data").
		Publish()

	assert.Nil(t, err)
}

func testMessagePublishError(t *testing.T) {
	defer nats.MockServer()()
	natsjs.Connect()
	natsjs.SetTimeout(0)

	err := natsjs.New().
		NewMessage().
		SetSubject("test", "subject").
		SetData("test data").
		Publish()

	assert.Error(t, err)
}

func testPublishNoConnection(t *testing.T) {
	natsjs.Close()

	assert.PanicsWithValue(
		t,
		"Connection is not established",
		func() {
			natsjs.New().
				NewMessage().
				SetSubject("test", "subject").
				SetData("test data").
				Publish()
		},
		"Publish should panic if connection is not set",
	)
}

func testPublishNoData(t *testing.T) {
	defer nats.MockServer()()
	natsjs.Connect()

	err := natsjs.New().
		NewMessage().
		SetSubject("test", "subject").
		Publish()

	assert.Error(t, err)
}
