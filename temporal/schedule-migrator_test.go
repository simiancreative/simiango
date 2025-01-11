package temporal_test

import (
	"testing"

	temporalMock "github.com/simiancreative/simiango/mocks/temporal"
	"github.com/simiancreative/simiango/temporal"
	"github.com/stretchr/testify/require"
)

var schedules temporal.ScheduleConfigs

func TestScheduleMigrate(t *testing.T) {
	server := temporalMock.MockService()
	temporalMock.SetHost(server)
	defer server.Stop()

	client := temporal.Connect()
	defer client.Close()

	// register simple workflow
	temporal.Register("simple")

	// set config to create a schedule
	setConfigData(temporal.ScheduleConfigs{
		{Workflow: "simple", ID: "1", Cron: "* * * * *", Input: "test"},
	})

	// run migrator
	err := temporal.Migrator(schedules)
	require.NoError(t, err)

	// check client.ScheduledConfigs() for the schedule
	configs, _ := client.ScheduledConfigs()
	require.Len(t, configs, len(schedules))

	// set empty config
	schedules = temporal.ScheduleConfigs{}

	// run migrator
	err = temporal.Migrator(schedules)

	// check client.ScheduledConfigs() for an empty list
	configs, _ = client.ScheduledConfigs()
	require.Len(t, configs, 0)
}

func setConfigData(configs temporal.ScheduleConfigs) {
	schedules = configs
}
