package gen

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aymerick/raymond"
	"github.com/spf13/cobra"

	"github.com/simiancreative/simiango/simian-go/app"
)

var kind string
var path string
var generators = map[string]Erator{}

var GenCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate project structure",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGen()
	},
}

func init() {
	GenCmd.Flags().StringVarP(&kind, "kind", "k", "service", "the generator")
	GenCmd.Flags().StringVarP(&path, "path", "p", ".", "the generator's working directory")

	app.RootCmd.AddCommand(GenCmd)
}

func runGen() error {
	generator := generators[kind]

	values, err := gatherValues(generator)
	if err != nil {
		return err
	}

	err = alertToWrite(generator, *values)
	if err != nil {
		return err
	}

	err = writeContent(generator, *values)
	if err != nil {
		return err
	}

	return nil
}

func writeContent(e Erator, v Values) error {
	for _, tpl := range e.Templates {
		if tpl.IF != "" && !v[tpl.IF].(bool) {
			continue
		}

		result, _ := raymond.Render(tpl.Content, v)
		content := []byte(result)
		err := ioutil.WriteFile(tpl.Path, content, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func alertToWrite(e Erator, v Values) error {
	var qs = []*survey.Question{{
		Name:     "ok",
		Validate: survey.Required,
		Prompt:   &survey.Confirm{Message: "Continue"},
	}}
	tpls := []string{}

	for _, tpl := range e.Templates {
		if tpl.IF != "" && !v[tpl.IF].(bool) {
			continue
		}

		tpls = append(tpls, tpl.Path)
	}

	answers := map[string]interface{}{}

	fmt.Println(fmt.Sprintf("Creating the following files in %s", path))
	fmt.Println(strings.Join(tpls, ", "))
	err := survey.Ask(qs, &answers)
	if err != nil {
		return err
	}

	val := answers["ok"].(bool)

	if !val {
		err = fmt.Errorf("stopped")
	}

	return err
}

func gatherValues(e Erator) (*Values, error) {
	var qs = []*survey.Question{}

	for _, rv := range e.RequiredVars {
		q := &survey.Question{
			Name:     rv.Name,
			Validate: survey.Required,
		}

		if rv.Type == "confirm" {
			q.Prompt = &survey.Confirm{Message: rv.Message}
		}

		if rv.Type != "confirm" {
			q.Prompt = &survey.Input{Message: rv.Message}
		}

		qs = append(qs, q)
	}

	answers := &map[string]interface{}{}

	err := survey.Ask(qs, answers)
	if err != nil {
		return nil, err
	}

	values := Values(*answers)

	return &values, nil
}

func Register(e Erator) {
	generators[e.Name] = e
}
