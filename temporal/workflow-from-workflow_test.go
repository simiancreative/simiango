package temporal_test

import (
	"testing"

	"github.com/simiancreative/simiango/workflow"

	"github.com/simiancreative/simiango/temporal"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
)

func TestRegisterActionsAsWorkfows(t *testing.T) {
	wfModel := workflow.Model{
		Name:        gofakeit.Word(),
		Description: gofakeit.Sentence(10),
		Actions: workflow.Actions{
			"action1": workflow.Action{
				Args:   workflow.ArgsList{},
				Runner: func(_ workflow.Args) (any, error) { return nil, nil },
			},
			"action2": workflow.Action{
				Args:   workflow.ArgsList{},
				Runner: func(_ workflow.Args) (any, error) { return nil, nil },
			},
		},
	}

	models := temporal.RegisterActionsAsWorkflows(wfModel)
	assert.Len(t, models, 2)
}

func TestRegisterWorkflowAsWorkflow(t *testing.T) {
	wfModel := workflow.Model{
		Name:        gofakeit.Word(),
		Description: gofakeit.Sentence(10),
		Actions: workflow.Actions{
			"action1": workflow.Action{
				Args:   workflow.ArgsList{},
				Runner: func(_ workflow.Args) (any, error) { return nil, nil },
			},
			"action2": workflow.Action{
				Args:   workflow.ArgsList{},
				Runner: func(_ workflow.Args) (any, error) { return nil, nil },
			},
		},
	}

	model := temporal.RegisterWorkflowAsWorkflow(wfModel)
	assert.NotNil(t, model)
}
