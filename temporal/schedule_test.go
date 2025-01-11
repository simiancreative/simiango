package temporal_test

import (
	"testing"

	temporalMock "github.com/simiancreative/simiango/mocks/temporal"
	"github.com/simiancreative/simiango/temporal"

	"github.com/stretchr/testify/require"
)

var assertions = map[string]func(client temporal.Client) func(t *testing.T){
	"TestSchedule_ID":       testScheduleID,
	"TestSchedule_Validate": testScheduleValidate,
	"TestSchedule_Remove":   testScheduleRemove,
	"TestSchedule_Upsert":   testScheduleUpsert,
}

func TestSchedules(t *testing.T) {
	server := temporalMock.MockService()
	temporalMock.SetHost(server)
	defer server.Stop()

	client := temporal.Connect()

	for name, test := range assertions {
		t.Run(name, test(client))
	}

	client.Close()
}

func testScheduleID(client temporal.Client) func(t *testing.T) {
	return func(t *testing.T) {
		schedule, err := client.Schedule("test-model", "123")

		require.NoError(t, err)

		expected := "test-model-123-scheduled"
		result := schedule.ID()

		require.Equal(t, expected, result)
	}
}

func testScheduleValidate(client temporal.Client) func(t *testing.T) {
	return func(t *testing.T) {
		schedule, err := client.Schedule("test-model", "123")
		require.NoError(t, err)

		err = schedule.Validate()
		require.Error(t, err)

		// sets cron
		schedule.SetCron("* * * * *")

		err = schedule.Validate()
		require.NoError(t, err)
	}
}

func testScheduleRemove(client temporal.Client) func(t *testing.T) {
	return func(t *testing.T) {
		schedule, err := client.Schedule("test-model", "123")
		require.NoError(t, err)

		err = schedule.SetCron("* * * * *").Upsert()
		require.NoError(t, err)

		err = schedule.Remove()
		require.NoError(t, err)
	}
}

func testScheduleUpsert(client temporal.Client) func(t *testing.T) {
	return func(t *testing.T) {
		schedule, err := client.Schedule("test-model", "123")
		require.NoError(t, err)

		err = schedule.SetCron("* * * * *").SetInput(`"input"`).Upsert()
		require.NoError(t, err)

		// should update the schedule
		err = schedule.SetCron("0 * * * *").Upsert()
		require.NoError(t, err)
	}
}
