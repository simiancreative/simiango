package config_test

import (
	"flag"
	"os"
	"testing"

	"github.com/simiancreative/simiango/config"
	"github.com/stretchr/testify/assert"
)

func TestWithFlag(t *testing.T) {
	// Backup the original args and restore them after the test
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	// Set the args for the test
	os.Args = []string{"cmd", "--dot-env=test"}
	flag.CommandLine = flag.NewFlagSet(os.Args[1], flag.ExitOnError)

	// Call the function
	config.WithFlag()

	// Check the environment variable
	assert.Equal(t, "test", os.Getenv("DOT_ENV"))
}

func TestWithFlagNoEnv(t *testing.T) {
	os.Unsetenv("DOT_ENV")

	// Backup the original args and restore them after the test
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	// Set the args for the test
	os.Args = []string{"cmd"}
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)

	// Call the function
	config.WithFlag()

	// Check the environment variable
	assert.Equal(t, "", os.Getenv("DOT_ENV"))
}

func TestEnable(t *testing.T) {
	// Set the environment variable
	os.Setenv("DOT_ENV", "test")

	// Call the function
	config.Enable()

	// Check the environment variable
	assert.NotEqual(t, "", os.Getenv("SOME_ENV_VAR"))
}

func TestEnableNoEnv(t *testing.T) {
	// Unset the environment variable
	os.Unsetenv("DOT_ENV")
	os.Unsetenv("SOME_ENV_VAR")

	// Call the function
	config.Enable()

	// Check the environment variable
	assert.Equal(t, "", os.Getenv("SOME_ENV_VAR"))
}
