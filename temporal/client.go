package temporal

import (
	"context"

	"github.com/simiancreative/simiango/errors"
	"github.com/simiancreative/simiango/logger"
	"go.temporal.io/sdk/client"
)

type Client struct {
	client.Client
}

func (c Client) HasModel(name string) error {
	_, err := findModel(name)
	return errors.Wrap(err, "Does not have model")
}

func (c Client) Start(name string) (stop func()) {
	model, err := findModel(name)
	if err != nil {
		panic(err)
	}

	model.client = c.Client

	// create worker
	return model.Listen()
}

func (c Client) ScheduledConfigs() (ScheduleConfigs, error) {
	instances := ScheduleConfigs{}

	list, err := c.
		Client.
		ScheduleClient().
		List(context.Background(), client.ScheduleListOptions{
			PageSize: 25,
		})

	if err != nil {
		return nil, errors.Wrap(err, "unable to list schedules")
	}

	for list.HasNext() {
		entry, err := list.Next()
		if err != nil {
			return nil, errors.Wrap(err, "unable to list schedules")
		}

		config := ScheduleConfig{}
		err = config.ParseMemo(entry.Memo)
		if err != nil {
			logger.Errorf("unable to create schedule config: %v", err)
			continue
		}

		instances.Add(config)
	}

	return instances, nil
}

func (c Client) Schedule(name, id string) (*Schedule, error) {
	model, err := findModel(name)
	if err != nil {
		return nil, errors.Wrap(err, "Schedule not found")
	}

	model.client = c.Client
	schedule := model.Schedule(id)

	return schedule, nil
}

func (c Client) Submit(name, input string) error {
	model, err := findModel(name)
	if err != nil {
		return errors.Wrap(err, "Submit failed")
	}

	model.client = c.Client

	return model.Submit(input)
}

func findModel(name string) (*Model, error) {
	model, ok := registered[name]

	if !ok {
		return nil, createModelNotFoundError(name)
	}

	return model, nil
}
