package temporal

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"time"

	sentrygo "github.com/getsentry/sentry-go"
	"github.com/simiancreative/simiango/errors"
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/monitoring/sentry"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	sdk "go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type Options = workflow.ActivityOptions
type RetryPolicy = sdk.RetryPolicy
type Workflow = func(input ActivityExecutor) (any, error)
type ActivityExecutor = func(activity any, input any) (any, error)

var DefaultRetryPolicy = &RetryPolicy{MaximumAttempts: 1}
var DefaultStartToCloseTimeout = time.Minute * 30

var DefaultOptions = Options{
	StartToCloseTimeout: DefaultStartToCloseTimeout,
	RetryPolicy:         DefaultRetryPolicy,
}

type Activity struct {
	Name string
	Func interface{}
}

// Model represents a Temporal Workflow.
type Model struct {
	client     client.Client
	Name       string
	Input      func() interface{}
	Workflow   Workflow
	Activities []Activity
	Options    Options
}

// SetOptions sets the ActivityOptions for the Model.
func (m *Model) SetOptions(o Options) *Model {
	m.Options = o

	return m
}

// SetInput sets the input object for the Model.
//
// The input object is used to unmarshal the input string when submitting a
// Workflow.
func (m *Model) SetInput(i func() interface{}) *Model {
	m.Input = i

	return m
}

// SetWorkflowFunc sets the Workflow function for the Model. When left unset, the
// Model will run a linear Workflow.
//
// The Workflow function is used to create a custom Workflow. it accepts a
// single argument, a function that executes an Activity and returns the output.
func (m *Model) SetWorkflowFunc(w Workflow) *Model {
	m.Workflow = w

	return m
}

func (m *Model) WorkflowID() string {
	return fmt.Sprintf("%s-%v", m.Name, time.Now().Unix())
}

func (m *Model) Schedule(id string) *Schedule {
	return &Schedule{id: id, model: m}
}

// Listen starts a Worker that listens on the Task Queue.
func (m Model) Listen() (stop func()) {
	// Create a Worker that listens on a Task Queue.
	w := worker.New(m.client, m.Name, worker.Options{})

	w.RegisterWorkflowWithOptions(m.runner, workflow.RegisterOptions{Name: m.Name})
	for _, a := range m.Activities {
		w.RegisterActivityWithOptions(a.Func, activity.RegisterOptions{Name: a.Name})
	}

	// Start listening to the Task Queue.
	err := w.Start()
	if err != nil {
		panic(errors.Wrap(err, "unable to start Worker"))
	}

	return w.Stop
}

// Submit starts a Workflow with the given input.
func (m *Model) Submit(inputStr string) error {
	input, err := m.prepareInput(inputStr)
	if err != nil {
		return errors.Wrap(err, "unable to prepare input")
	}

	options := client.StartWorkflowOptions{
		ID:        m.WorkflowID(),
		TaskQueue: m.Name,
	}

	we, err := m.client.ExecuteWorkflow(context.Background(), options, m.runner, input)
	if err != nil {
		return errors.Wrap(err, "unable to start the Workflow")
	}

	var result interface{}

	err = we.Get(context.Background(), &result)
	if err != nil {
		return errors.Wrap(err, "workflow execution failed")
	}

	logger.Debugf("WorkflowID: %s RunID: %s", we.GetID(), we.GetRunID())

	return nil
}

func (m Model) runner(ctx workflow.Context, input interface{}) (interface{}, error) {
	sentry.SetContext("workflow", map[string]interface{}{
		"name":      m.Name,
		"stepName":  "temporal-runner",
		"arguments": input,
	})

	options := m.Options
	ctx = workflow.WithActivityOptions(ctx, options)

	runner := linearWorkflow

	if m.Workflow != nil {
		runner = customWorkflow
	}

	result, err := runner(m, ctx, input)
	if err != nil {
		sentrygo.CaptureException(errors.Unwrap(err))
	}

	return result, err
}

func customWorkflow(m Model, ctx workflow.Context, originalInput interface{}) (interface{}, error) {
	execute := func(activity interface{}, input interface{}) (interface{}, error) {
		var name string
		activityType := reflect.TypeOf(activity)
		if activityType.Kind() == reflect.Func {
			funcPtr := reflect.ValueOf(activity).Pointer()
			funcName := runtime.FuncForPC(funcPtr).Name()
			name = funcName
		}

		var output interface{}
		if input == nil {
			input = originalInput
		}

		logger.Debugf("Executing activity func: %v", name)

		sentry.SetContext("workflow", map[string]interface{}{
			"name":      m.Name,
			"stepName":  name,
			"arguments": input,
		})

		err := workflow.
			ExecuteActivity(ctx, activity, input).
			Get(ctx, &output)

		logger.Debugf("Activity output: %v, %v", name, output)

		if err != nil {
			return nil, errors.Wrap(err, "activity failed")
		}

		logger.Debugf("Activity succeeded: %s", name)
		return output, nil
	}

	return m.Workflow(execute)
}

func linearWorkflow(m Model, ctx workflow.Context, input interface{}) (interface{}, error) {
	output := input

	logger.Debugf("Running linear workflow: %s", m.Name)

	for _, activity := range m.Activities {
		logger.Debugf("Executing activity: %s", activity.Name)

		sentry.SetContext("workflow", map[string]interface{}{
			"name":      m.Name,
			"stepName":  activity.Name,
			"arguments": output,
		})

		err := workflow.
			ExecuteActivity(ctx, activity.Func, output).
			Get(ctx, &output)

		logger.Debugf("Activity output: %v", output)

		if err != nil {
			logger.Debugf("Activity failed: %s, %#v", activity.Name, err)
			return nil, errors.Wrap(err, "activity failed")
		}
	}

	return output, nil
}

func (m *Model) prepareInput(inputStr string) (interface{}, error) {
	var input interface{}

	if m.Input != nil && len(inputStr) > 0 {
		input = m.Input()
		err := json.Unmarshal([]byte(inputStr), input)
		if err != nil {
			return nil, errors.Wrap(err, "unable to unmarshal input")
		}
	}

	return input, nil
}
