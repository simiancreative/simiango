package temporal

import (
	"encoding/json"
	"fmt"

	"github.com/simiancreative/simiango/errors"

	"go.temporal.io/api/common/v1"
)

// ScheduleConfig represents a configuration for a scheduled workflow
type ScheduleConfig struct {
	Comment  string
	Workflow string
	ID       string
	Cron     string
	Input    string
}

func (s ScheduleConfig) Event() string {
	event := map[string]string{"event": s.Input}
	eventBytes, _ := json.Marshal(event)

	return string(eventBytes)
}

func (s ScheduleConfig) Name() string {
	workflow, err := findModel(s.Workflow)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s-%s", workflow.Name, s.ID)
}

func (s *ScheduleConfig) ParseMemo(memo *common.Memo) error {
	fields := memo.GetFields()

	data, ok := fields["workflow"]
	if !ok {
		return errors.New("unable to find workflow name in schedule")
	}

	workflowByte := data.GetData()
	workflow := ""
	err := json.Unmarshal(workflowByte, &workflow)
	if err != nil {
		return errors.New("unable to unmarshal workflow name")
	}

	data, ok = fields["id"]
	if !ok {
		return errors.New("unable to find id in schedule")
	}

	idByte := data.GetData()
	id := ""
	err = json.Unmarshal(idByte, &id)
	if err != nil {
		return errors.New("unable to unmarshal id")
	}

	*s = ScheduleConfig{
		Workflow: workflow,
		ID:       id,
	}

	return nil
}

// ScheduleConfigs is a collection of ScheduleConfig
type ScheduleConfigs []ScheduleConfig

func (configs *ScheduleConfigs) Add(config ScheduleConfig) {
	*configs = append(*configs, config)
}

func (configs ScheduleConfigs) Diff(comparable ScheduleConfigs) ScheduleConfigs {
	remainder := &ScheduleConfigs{}

	for _, config := range configs {
		if comparable.Contains(config) {
			continue
		}

		remainder.Add(config)
	}

	return *remainder
}

func (configs ScheduleConfigs) Contains(config ScheduleConfig) bool {
	for _, i := range configs {
		if i.Name() == config.Name() {
			return true
		}
	}

	return false
}
