package uuid

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/simiancreative/simiango/meta"
	metacli "github.com/simiancreative/simiango/simian-go/app/meta"
)

var key string
var tokenStr string

var cmd = &cobra.Command{
	Use:   "uuid",
	Short: "generate a uuid",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	metacli.Cmd.AddCommand(cmd)
}

func run() error {
	var style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF")).
		Background(lipgloss.Color("#158978")).
		Bold(true).
		MarginTop(1).
		PaddingTop(1).
		PaddingRight(4).
		PaddingBottom(1).
		PaddingLeft(4)

	fmt.Println(style.Render(string(meta.Id())))

	return nil
}
