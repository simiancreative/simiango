package temporal_test

import (
	"testing"

	"github.com/simiancreative/simiango/errors"
	"github.com/simiancreative/simiango/mocks/temporal"
	mock "github.com/simiancreative/simiango/mocks/temporal"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListen_FailedStart(t *testing.T) {
	id := uuid.New().String()

	// create the mock service
	s := mock.MockService()

	// register a model
	temporal.Register(id)

	// set the host and start the worker
	client := temporal.Connect(s.FrontendHostPort())
	defer client.Close()

	// stop the mock service
	s.Stop()

	assert.Panics(t, func() {
		client.Start(id)
	}, "Listen should panic if the client is not set")
}

func TestSubmit_FailedMarshal(t *testing.T) {
	id := uuid.New().String()

	// register a model
	err := temporal.
		Register(id).
		SetInput(func() any { return &map[string]any{} }).
		Submit("{")

	assert.Error(t, err, "Submit should return an error if the input cannot be marshaled")
}

func TestSubmit_FailedExecute(t *testing.T) {
	id := uuid.New().String()

	// create the mock service
	s := mock.MockService()

	// register a model
	temporal.Register(id)

	// set the host and start the worker
	client := temporal.Connect(s.FrontendHostPort())
	defer client.Close()

	// stop the mock service
	s.Stop()

	err := client.Submit(id, "{}")

	assert.Error(t, err, "Submit should return an error if the client is not set")
}

func TestSubmit_WorkflowFailure(t *testing.T) {
	id := uuid.New().String()

	// create the mock service
	s := mock.MockService()
	defer s.Stop()

	// register a model
	temporal.
		Register(id).
		SetWorkflowFunc(func(_ temporal.ActivityExecutor) (any, error) {
			return nil, errors.New("workflow failure")
		})

	// set the host and start the worker
	client := temporal.Connect(s.FrontendHostPort())
	defer client.Close()

	stop := client.Start(id)
	defer stop()

	err := client.Submit(id, "{}")

	assert.Error(t, err, "Submit should return an error")
}

func TestSubmit_ActivityFailure(t *testing.T) {
	id := uuid.New().String()

	// create the mock service
	s := mock.MockService()
	defer s.Stop()

	activity := func(_ any, _ any) (any, error) {
		return nil, errors.New("activity failure")
	}

	// register a model
	temporal.Register(id, temporal.Activity{Name: "activity", Func: activity}).
		SetOptions(temporal.DefaultOptions)

	// set the host and start the worker
	client := temporal.Connect(s.FrontendHostPort())
	defer client.Close()

	stop := client.Start(id)
	defer stop()

	err := client.Submit(id, "{}")

	assert.Error(t, err, "Submit should return an error")
}

func TestSubmit_CustomWorkflow(t *testing.T) {
	id := uuid.New().String()

	// create the mock service
	s := mock.MockService()
	defer s.Stop()

	activity := func(_ any, _ any) (any, error) {
		return 1, nil
	}

	wf := func(exec temporal.ActivityExecutor) (any, error) {
		return exec(activity, nil)
	}

	// register a model
	temporal.Register(id, temporal.Activity{Name: "activity", Func: activity}).
		SetWorkflowFunc(wf).
		SetOptions(temporal.DefaultOptions)

	// set the host and start the worker
	client := temporal.Connect(s.FrontendHostPort())
	defer client.Close()

	stop := client.Start(id)
	defer stop()

	err := client.Submit(id, "{}")
	assert.NoError(t, err, "Submit should not return an error")
}

func TestSubmit_CustomWorkflowActivityFailure(t *testing.T) {
	id := uuid.New().String()

	// create the mock service
	s := mock.MockService()
	defer s.Stop()

	activity := func(_ any, _ any) (any, error) {
		return nil, errors.New("activity failure")
	}

	wf := func(exec temporal.ActivityExecutor) (any, error) {
		return exec(activity, nil)
	}

	// register a model
	temporal.Register(id, temporal.Activity{Name: "activity", Func: activity}).
		SetWorkflowFunc(wf).
		SetOptions(temporal.DefaultOptions)

	// set the host and start the worker
	client := temporal.Connect(s.FrontendHostPort())
	defer client.Close()

	stop := client.Start(id)
	defer stop()

	err := client.Submit(id, "{}")
	assert.Error(t, err, "Submit should return an error")
}
