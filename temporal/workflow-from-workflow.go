package temporal

import (
	"fmt"

	"github.com/simiancreative/simiango/workflow"
)

func RegisterActionsAsWorkflows(model workflow.Model) map[string]*Model {
	models := map[string]*Model{}

	for name, action := range model.Actions {
		name = fmt.Sprintf("%s-%s", model.Name, name)

		model := Register(name, Activity{Name: name, Func: action.Runner}).
			SetInput(func() interface{} {
				return &workflow.Args{}
			}).
			SetOptions(DefaultOptions)

		models[name] = model
	}

	return models
}

func RegisterWorkflowAsWorkflow(model workflow.Model) *Model {
	activities := []Activity{}
	for name, action := range model.Actions {
		activities = append(
			activities,
			Activity{Name: name, Func: action.Runner},
		)
	}

	return Register(model.Name, activities...).
		SetInput(func() interface{} {
			return &workflow.Args{}
		}).
		SetOptions(DefaultOptions)
}
