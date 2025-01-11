package temporal

import (
	"context"
	"strings"

	"github.com/simiancreative/simiango/errors"
	"github.com/simiancreative/simiango/logger"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

type Schedule struct {
	model  *Model
	id     string
	handle client.ScheduleHandle
	cron   string
	input  string
}

func (s *Schedule) SetCron(cron string) *Schedule {
	s.cron = cron

	return s
}

func (s *Schedule) SetInput(input string) *Schedule {
	s.input = input

	return s
}

func (s *Schedule) ID() string {
	m := s.model
	parts := []string{m.Name, s.id, "scheduled"}
	return strings.Join(parts, "-")
}

func (s *Schedule) Validate() error {
	if s.cron == "" {
		return errors.New("cron is required")
	}

	return nil
}

func (s *Schedule) Remove() error {
	// dont check for cron here, just remove the schedule
	m := s.model
	client := m.client.ScheduleClient()

	handle := client.GetHandle(context.Background(), s.ID())

	if handle == nil {
		return nil
	}

	logger.Warnf("Removing schedule: %v", s.ID())

	err := handle.Delete(context.Background())

	if err != nil {
		return errors.Wrap(err, "unable to remove the Workflow schedule")
	}

	return nil
}

func (s *Schedule) Upsert() error {
	m := s.model

	// only block creation if cron is not set
	if err := s.Validate(); err != nil {
		return errors.Wrap(err, "upsert failed")
	}

	client := m.client.ScheduleClient()

	s.handle = client.GetHandle(context.Background(), s.ID())
	_, err := s.handle.Describe(context.Background())

	if err == nil {
		return s.updateSchedule()
	}

	return s.createSchedule()
}

func (s *Schedule) updateSchedule() error {
	logger.Warnf("Updating schedule: %v", s.ID())

	options, err := s.scheduleOptions()
	if err != nil {
		return errors.Wrap(err, "unable to prepare schedule options")
	}

	doUpdate := func(input client.ScheduleUpdateInput) (*client.ScheduleUpdate, error) {
		input.Description.Schedule.Spec = &options.Spec
		input.Description.Schedule.Action = options.Action
		input.Description.Schedule.Policy.Overlap = options.Overlap

		return &client.ScheduleUpdate{
			Schedule: &input.Description.Schedule,
		}, nil
	}

	err = s.handle.Update(context.Background(), client.ScheduleUpdateOptions{
		DoUpdate: doUpdate,
	})

	if err != nil {
		return errors.Wrap(err, "unable to update the Workflow schedule")
	}

	return nil
}

func (s *Schedule) createSchedule() error {
	m := s.model
	logger.Warnf("Creating schedule: %v", s.ID())

	options, err := s.scheduleOptions()
	if err != nil {
		return errors.Wrap(err, "unable to prepare schedule options")
	}

	scheduleClient := m.client.ScheduleClient()

	_, err = scheduleClient.Create(context.Background(), *options)

	if err != nil {
		return errors.Wrap(err, "unable to create the Workflow schedule")
	}

	return nil
}

func (s *Schedule) scheduleOptions() (*client.ScheduleOptions, error) {
	m := s.model

	input, err := m.prepareInput(s.input)
	if err != nil {
		return nil, errors.Wrap(err, "unable to prepare input")
	}

	return &client.ScheduleOptions{
		ID: s.ID(),
		Spec: client.ScheduleSpec{
			CronExpressions: []string{s.cron},
		},
		Action: &client.ScheduleWorkflowAction{
			ID:                  m.Name,
			Workflow:            m.runner,
			Args:                []interface{}{input},
			RetryPolicy:         m.Options.RetryPolicy,
			WorkflowTaskTimeout: m.Options.StartToCloseTimeout,
			TaskQueue:           m.Name,
		},
		Overlap: enums.SCHEDULE_OVERLAP_POLICY_SKIP,
		Memo: map[string]interface{}{
			"workflow": m.Name,
			"id":       s.id,
		},
	}, nil
}
