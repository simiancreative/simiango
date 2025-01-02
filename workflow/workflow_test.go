package workflow_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/simiancreative/simiango/cli"
	"github.com/simiancreative/simiango/workflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	testAction := workflow.Action{
		Args: workflow.ArgsList{
			workflow.ArgsItem{"test", "test Description"},
		},
		Runner: func(_ workflow.Args) (interface{}, error) {
			return "testisawesome", nil
		},
	}

	model := workflow.Model{
		Name:        "exec",
		Description: "test workflow",
		Actions: workflow.Actions{
			"runner": testAction,
		},
	}

	root := cli.New("workflows")

	model.RegisterAsCommand(root.Cmd)

	b := bytes.NewBufferString("")

	root.Cmd.SetOut(b)
	root.Cmd.SetArgs([]string{"exec", "runner", "23"})

	root.Execute()

	out, err := io.ReadAll(b)

	assert.NoError(t, err)
	assert.Equal(t, "\"testisawesome\"", string(out))
}

func TestValidFlags(t *testing.T) {
	testAction := workflow.Action{
		Args: workflow.ArgsList{
			workflow.ArgsItem{"test", "test Description"},
		},
		Runner: func(_ workflow.Args) (interface{}, error) {
			return "testisawesome", nil
		},
	}

	model := workflow.Model{
		Name:        "valid-flags",
		Description: "test workflow",
		Actions: workflow.Actions{
			"runner": testAction,
		},
	}

	root := cli.New("workflows")

	model.RegisterAsCommand(root.Cmd)

	b := bytes.NewBufferString("")

	root.Cmd.SetOut(b)
	root.Cmd.SetArgs([]string{"valid-flags", "runner", "--test", "23"})

	root.Execute()

	out, err := io.ReadAll(b)

	assert.NoError(t, err)
	assert.Equal(t, "\"testisawesome\"", string(out))
}

func TestFailedArgs(t *testing.T) {
	testAction := workflow.Action{
		Args: workflow.ArgsList{
			workflow.ArgsItem{"test", "test Description"},
		},
		Runner: func(_ workflow.Args) (interface{}, error) {
			return "testisawesome", nil
		},
	}

	model := workflow.Model{
		Name:        "invalid-args-tester",
		Description: "test workflow",
		Actions: workflow.Actions{
			"invalid-args": testAction,
		},
	}

	root := cli.New("workflows")

	model.RegisterAsCommand(root.Cmd)

	root.Cmd.SetArgs([]string{"invalid-args-tester", "invalid-args"})

	err := root.Cmd.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "Args/Flags not found for (invalid-args-tester:invalid-args)")
}

func TestSomeValidArgs(t *testing.T) {
	testAction := workflow.Action{
		Args: workflow.ArgsList{
			workflow.ArgsItem{"test", "test Description"},
			workflow.ArgsItem{"test-two", "test Description"},
		},
		Runner: func(_ workflow.Args) (interface{}, error) {
			return "testisawesome", nil
		},
	}

	model := workflow.Model{
		Name:        "some-valid-args",
		Description: "test workflow",
		Actions: workflow.Actions{
			"runner": testAction,
		},
	}

	root := cli.New("workflows")

	model.RegisterAsCommand(root.Cmd)

	root.Cmd.SetArgs([]string{"some-valid-args", "runner", "23"})

	err := root.Cmd.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "Args/Flags not found for (some-valid-args:runner)")
}

func TestHandleError(t *testing.T) {
	testAction := workflow.Action{
		Args: workflow.ArgsList{
			workflow.ArgsItem{"test", "test Description"},
		},
		Runner: func(_ workflow.Args) (interface{}, error) {
			return nil, errors.New("test error")
		},
	}

	model := workflow.Model{
		Name:        "handle-error",
		Description: "test workflow",
		Actions: workflow.Actions{
			"runner": testAction,
		},
	}

	root := cli.New("workflows")

	model.RegisterAsCommand(root.Cmd)

	b := bytes.NewBufferString("")

	root.Cmd.SetOut(b)
	root.Cmd.SetArgs([]string{"handle-error", "runner", "23"})

	err := root.Cmd.Execute()

	assert.Error(t, err)
}

func TestMixOfFlagsAndPositional(t *testing.T) {
	testAction := workflow.Action{
		Args: workflow.ArgsList{
			workflow.ArgsItem{"test", "test Description"},
			workflow.ArgsItem{"test-two", "test Description"},
		},
		Runner: func(args workflow.Args) (interface{}, error) {
			return args, nil
		},
	}

	model := workflow.Model{
		Name:        "mix-args",
		Description: "test workflow",
		Actions: workflow.Actions{
			"runner": testAction,
		},
	}

	b := bytes.NewBufferString("")

	root := cli.New("workflows")

	model.RegisterAsCommand(root.Cmd)

	root.Cmd.SetOut(b)
	root.Cmd.SetArgs([]string{"mix-args", "runner", "--test-two", "23", "24"})

	err := root.Cmd.Execute()
	require.NoError(t, err)

	out, err := io.ReadAll(b)

	require.NoError(t, err)
	require.Contains(t, string(out), `"test":"24","test-two":"23"`)
}

func TestItemFromJSON(t *testing.T) {
	args := workflow.Args{
		"key1": `{"subkey1":"value1"}`,
	}

	result := args.ItemFromJSON("key1")
	expected := `{"subkey1":"value1"}`

	assert.Equal(t, expected, result)
}

func TestItemFromJSONWithInvalidKey(t *testing.T) {
	args := workflow.Args{
		"key1": `{"subkey1":"value1"}`,
	}

	result := args.ItemFromJSON("key2")

	assert.Empty(t, result)
}

func TestUnmarshalKey(t *testing.T) {
	args := workflow.Args{
		"key1": `{"subkey1":"value1"}`,
	}

	var dest map[string]string
	err := args.UnmarshalKey("key1", &dest)

	assert.NoError(t, err)
	assert.Equal(t, "value1", dest["subkey1"])
}

func TestUnmarshalKeyWithInvalidKey(t *testing.T) {
	args := workflow.Args{
		"key1": `{"subkey1":"value1"}`,
	}

	var dest map[string]string
	err := args.UnmarshalKey("key2", &dest)

	assert.Error(t, err)
}
