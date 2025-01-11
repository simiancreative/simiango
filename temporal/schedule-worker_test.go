package temporal_test

import (
	"os"
	"testing"
	"time"

	"github.com/simiancreative/simiango/config"
	temporalmock "github.com/simiancreative/simiango/mocks/temporal"
	"github.com/simiancreative/simiango/temporal"
	"github.com/simiancreative/simiango/workflow"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	name   string
	params workflow.Args
	method func(t *testing.T, params workflow.Args)
}{
	{
		name: "valid process",
		params: workflow.Args{
			"workflow": uuid.New().String(),
		},
		method: func(t *testing.T, params workflow.Args) {
			temporal.Register(params["workflow"]).SetOptions(temporal.DefaultOptions)
			result, err := temporal.ScheduleWorker(params["workflow"])

			assert.NoError(t, err)
			assert.Nil(t, result)
		},
	},
	{
		name: "invalid workflow",
		params: workflow.Args{
			"workflow": "invalid",
		},
		method: func(t *testing.T, params workflow.Args) {
			result, err := temporal.ScheduleWorker(params["workflow"])

			assert.Error(t, err)
			assert.Nil(t, result)
		},
	},
	{
		name: "invalid stream",
		params: workflow.Args{
			"workflow": "invalid",
		},
		method: func(t *testing.T, params workflow.Args) {
			result, err := temporal.ScheduleWorker(params["workflow"])

			assert.Error(t, err)
			assert.Nil(t, result)
		},
	},
}

func TestWorker(t *testing.T) {
	s := temporalmock.MockService()
	temporalmock.SetHost(s)
	defer s.Stop()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.SetupTest()

			proc, err := os.FindProcess(os.Getpid())
			if err != nil {
				t.Fatal(err)
			}

			go func() {
				time.Sleep(1 * time.Second)
				proc.Signal(os.Interrupt)
			}()

			tt.method(t, tt.params)
		})
	}
}
