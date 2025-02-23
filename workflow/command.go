package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/simiancreative/simiango/cli"
	"github.com/simiancreative/simiango/errors"
	"github.com/simiancreative/simiango/monitoring/sentry"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func (w Model) RegisterAsCommand(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   w.Name,
		Short: w.Description,
	}

	for name, action := range w.Actions {
		step := &cobra.Command{
			Use:  name,
			RunE: w.command(name, action),
		}

		for _, arg := range action.Args {
			step.Flags().String(arg[0], "", arg[1])
		}

		cmd.AddCommand(step)
	}

	rootCmd.AddCommand(cmd)
}

func (w Model) command(
	stepName string,
	action Action,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		sentry.SetContext("workflow", map[string]interface{}{
			"name":      w.Name,
			"stepName":  stepName,
			"arguments": args,
		})

		defer cli.Recover()

		if err := checkArgs(cmd.Flags(), args, w.Name, stepName, action.Args); err != nil {
			return errors.Wrap(err, "failed to check args")
		}

		val, err := action.Runner(buildArgs(cmd.Flags(), args, action.Args))
		if err != nil {
			return errors.Wrap(err, "failed to run action")
		}

		jsonBytes, _ := json.Marshal(val)

		fmt.Fprint(cmd.OutOrStdout(), string(jsonBytes))
		return nil
	}
}

func checkArgs(
	flags *pflag.FlagSet,
	args []string,
	name, stepName string,
	actionArgs ArgsList,
) error {
	if (flags.NFlag() + len(args)) >= len(actionArgs) {
		return nil
	}

	list := "\n"
	for i, val := range actionArgs {
		status := "❌"
		if len(args) >= (i+1) || flags.Changed(val[0]) {
			status = "✅"
		}

		list = fmt.Sprintf("%v%v %v : %v\n", list, status, val[0], val[1])
	}

	return fmt.Errorf(
		"Args/Flags not found for (%v:%v) required args are:\n %v",
		name, stepName, list,
	)
}

func buildArgs(flags *pflag.FlagSet, args []string, actionArgs ArgsList) Args {
	mapped := Args{}

	for i, val := range actionArgs {
		if flags.Changed(val[0]) {
			mapped[val[0]] = flags.Lookup(val[0]).Value.String()
			continue
		}

		mapped[val[0]] = args[i]
	}

	return mapped
}
