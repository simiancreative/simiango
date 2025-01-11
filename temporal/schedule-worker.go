package temporal

import (
	"github.com/simiancreative/simiango/errors"
	"github.com/simiancreative/simiango/sig"
)

func ScheduleWorker(workflow string) (interface{}, error) {
	// connect to Temporal
	client := Connect()

	if err := client.HasModel(workflow); err != nil {
		return nil, errors.Wrap(err, "unable to find workflow")
	}

	// Start the worker to listen for submitted workflows
	workerStop := Connect().
		Start(workflow)

	// add cleanup functions
	_, exit := sig.
		New().
		AddCleanup(workerStop).
		AddCleanup(client.Close).
		Catch()

	// wait for exit signal
	<-exit.Done()

	return nil, nil
}
