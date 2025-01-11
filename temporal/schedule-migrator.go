package temporal

import (
	"github.com/simiancreative/simiango/logger"
)

func Migrator(configuration ScheduleConfigs) error {
	client := Connect()

	// get all temporal schedules
	temporalSchedules, err := client.ScheduledConfigs()
	if err != nil {
		return err
	}

	// find the difference between the two sets
	needsDeletion := temporalSchedules.Diff(configuration)

	// delete all schedules that are not in the registered schedules
	for _, inst := range needsDeletion {
		schedule, err := client.Schedule(inst.Workflow, inst.ID)
		if err != nil {
			logger.Errorf("unable to find schedule %s: %v", inst.Name(), err)
		}

		err = schedule.Remove()
		if err != nil {
			logger.Errorf("unable to delete schedule %s: %v", inst.Name(), err)
		}
	}

	// upsert all schedules that are in the registered schedules
	for _, inst := range configuration {
		schedule, err := client.Schedule(inst.Workflow, inst.ID)
		if err != nil {
			logger.Errorf("unable to find schedule %s: %v", inst.Name(), err)
		}

		err = schedule.
			SetCron(inst.Cron).
			SetInput(inst.Event()).
			Upsert()

		if err != nil {
			logger.Errorf("unable to upsert schedule %s: %v", inst.Name(), err)
		}
	}

	return nil
}
