package temporal_test

import (
	"testing"

	"github.com/simiancreative/simiango/temporal"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
	"go.temporal.io/api/common/v1"
)

func init() {
	temporal.
		Register("test-model", temporal.Activity{Name: "test-model"}).
		SetInput(nil)
}

func TestScheduleConfig_Event(t *testing.T) {
	instance := temporal.ScheduleConfig{Input: "test-input"}
	expected := `{"event":"test-input"}`

	result := instance.Event()

	require.Equal(t, expected, result)
}

func TestScheduleConfig_Name(t *testing.T) {
	instance := temporal.ScheduleConfig{Workflow: "test-model", ID: "123"}
	expected := "test-model-123"

	result := instance.Name()
	require.Equal(t, expected, result)

	// Test error handling
	instance = temporal.ScheduleConfig{Workflow: "invalid-workflow", ID: "123"}
	assert.Panics(t, func() { instance.Name() })
}

func TestScheduleConfigs_Add(t *testing.T) {
	instances := &temporal.ScheduleConfigs{}
	instance := temporal.ScheduleConfig{ID: "123"}

	instances.Add(instance)

	if len(*instances) != 1 {
		t.Errorf("expected 1 instance, got %d", len(*instances))
	}

	if (*instances)[0].ID != "123" {
		t.Errorf("expected instance ID 123, got %s", (*instances)[0].ID)
	}
}

func TestScheduleConfigs_Diff(t *testing.T) {
	instance1 := temporal.ScheduleConfig{Workflow: "test-model", ID: "123"}
	instance2 := temporal.ScheduleConfig{Workflow: "test-model", ID: "456"}
	instances := temporal.ScheduleConfigs{instance1, instance2}
	comparable := temporal.ScheduleConfigs{instance1}

	result := instances.Diff(comparable)

	require.Len(t, result, 1)
	require.Equal(t, instance2, result[0])
}

func TestScheduleConfigs_Contains(t *testing.T) {
	instance1 := temporal.ScheduleConfig{Workflow: "test-model", ID: "123"}
	instance2 := temporal.ScheduleConfig{Workflow: "test-model", ID: "456"}
	instances := temporal.ScheduleConfigs{instance1}

	require.True(t, instances.Contains(instance1))
	require.False(t, instances.Contains(instance2))
}

func TestScheduleConfig_ParseMemo(t *testing.T) {
	memo := &common.Memo{
		Fields: map[string]*common.Payload{
			"workflow": {Data: []byte(`"test-model"`)},
			"id":       {Data: []byte(`"123"`)},
		},
	}

	instance := temporal.ScheduleConfig{}
	err := instance.ParseMemo(memo)

	require.NoError(t, err)
	require.Equal(t, "test-model", instance.Workflow)
	require.Equal(t, "123", instance.ID)
}

func TestScheduleConfig_ParseMemo_Error(t *testing.T) {
	memos := []*common.Memo{
		{Fields: map[string]*common.Payload{}},
		{Fields: map[string]*common.Payload{"workflow": {Data: []byte(`123`)}}},
		{Fields: map[string]*common.Payload{"workflow": {Data: []byte(`"123"`)}}},
		{
			Fields: map[string]*common.Payload{
				"workflow": {Data: []byte(`"test-model"`)},
				"id":       {Data: []byte(`123`)},
			},
		},
	}

	for _, m := range memos {
		instance := temporal.ScheduleConfig{}
		err := instance.ParseMemo(m)

		require.Error(t, err)
	}
}
