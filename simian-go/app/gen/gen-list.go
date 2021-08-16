package gen

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	// "github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/simiancreative/simiango/simian-go/app"
)

var GenListCmd = &cobra.Command{
	Use:   "gen-list",
	Short: "list available generators",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGenList()
	},
}

func init() {
	app.RootCmd.AddCommand(GenListCmd)
}

func runGenList() error {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"Kind", "Desc"})

	for index, generator := range generators {
		t.AppendRow(table.Row{index, generator.Desc})
	}

	fmt.Println("")
	fmt.Println(t.Render())

	return nil
}
