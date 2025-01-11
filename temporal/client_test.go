package temporal_test

import (
	"context"
	"testing"

	mock "github.com/simiancreative/simiango/mocks/temporal"
	"github.com/simiancreative/simiango/temporal"

	"github.com/tj/assert"
	"go.temporal.io/api/enums/v1"
	temporalClient "go.temporal.io/sdk/client"
)

func TestStart(t *testing.T) {
	s := mock.MockService()
	defer s.Stop()

	temporal.Register("startTest")

	client := temporal.Connect(s.FrontendHostPort())
	defer client.Close()

	assert.NotPanics(t, func() {
		client.Start("startTest")
	})

	assert.Panics(t, func() {
		client.Start("nonExistentModel")
	})
}

func TestSubmit(t *testing.T) {
	service := mock.MockService()
	defer service.Stop()

	temporal.Register("unitTestModel", temporal.Activity{
		Name: "unitTestActivity",
		Func: func(_ any, _ any) (any, error) {
			return nil, nil
		},
	}).SetOptions(temporal.DefaultOptions)

	client := temporal.Connect(service.FrontendHostPort())
	defer client.Close()

	defer client.Start("unitTestModel")()

	err := client.Submit("unitTestModel", `{}`)
	assert.NoError(t, err)
}

func TestSubmit_NoModel(t *testing.T) {
	client := temporal.Client{}
	err := client.Submit("nonExistentModel", "input")

	assert.Error(t, err)
}

func TestClose(t *testing.T) {
	s := mock.MockService()
	defer s.Stop()

	client := temporal.Connect(s.FrontendHostPort())
	defer client.Close()

	assert.NotPanics(t, func() {
		client.Close()
	}, "Close should not panic")
}

func TestSchedule(t *testing.T) {
	s := mock.MockService()
	defer s.Stop()

	temporal.Register("scheduleTest")

	client := temporal.Connect(s.FrontendHostPort())
	defer client.Close()

	_, err := client.Schedule("scheduleTest", "123")
	assert.NoError(t, err)

	_, err = client.Schedule("nonExistentModel", "123")
	assert.Error(t, err)
}

func TestScheduledConfigs(t *testing.T) {
	s := mock.MockService()
	defer s.Stop()

	temporal.Register("scheduledConfigsTest")

	client := temporal.Connect(s.FrontendHostPort())

	// createa valid schedule
	schedule, err := client.Schedule("scheduledConfigsTest", "123")
	assert.NoError(t, err)

	schedule.SetCron("* * * * *").SetInput("input").Upsert()

	// get the scheduled configs
	_, err = client.ScheduledConfigs()
	assert.NoError(t, err)

	// create a schedule without a memo
	scheduleClient := client.ScheduleClient()

	_, err = scheduleClient.Create(context.Background(), temporalClient.ScheduleOptions{
		ID: "test",
		Spec: temporalClient.ScheduleSpec{
			CronExpressions: []string{"* * * * *"},
		},
		Action: &temporalClient.ScheduleWorkflowAction{
			ID:                  "test",
			Workflow:            func() {},
			Args:                []interface{}{""},
			RetryPolicy:         temporal.DefaultRetryPolicy,
			WorkflowTaskTimeout: temporal.DefaultStartToCloseTimeout,
			TaskQueue:           "test",
		},
		Overlap: enums.SCHEDULE_OVERLAP_POLICY_SKIP,
	})

	// cant parse but continues
	_, err = client.ScheduledConfigs()
	assert.NoError(t, err)

	client.Close()

	// cant connect so returns an error
	_, err = client.ScheduledConfigs()
	assert.Error(t, err)
}

func TestClient_HasModel(t *testing.T) {
	s := mock.MockService()
	defer s.Stop()

	temporal.Register("hasModelTest")

	client := temporal.Connect(s.FrontendHostPort())
	defer client.Close()

	assert.NoError(t, client.HasModel("hasModelTest"))
	assert.Error(t, client.HasModel("nonExistentModel"))
}
