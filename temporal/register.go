package temporal

import (
	"github.com/simiancreative/simiango/errors"
	"github.com/simiancreative/simiango/workflow"
)

var registered = map[string]*Model{}

func Register(name string, activities ...Activity) *Model {
	model := &Model{
		Name:       name,
		Activities: activities,
		Options:    DefaultOptions,
		Input:      func() interface{} { return &workflow.Args{} },
	}

	if _, ok := registered[name]; ok {
		panic(errors.New("model already registered: %v", name))
	}

	registered[name] = model

	return model
}
