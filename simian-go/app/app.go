package app

import (
	"github.com/simiancreative/simiango/cli"
)

// RootCmd represents the base command when called without any subcommands
var Root = cli.New("simian-go")

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
var Execute = Root.Cmd.Execute
