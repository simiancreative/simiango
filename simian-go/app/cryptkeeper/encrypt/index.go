package encrypt

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/charmbracelet/lipgloss"
	"github.com/simiancreative/simiango/cryptkeeper"
	cryptkeepercli "github.com/simiancreative/simiango/simian-go/app/cryptkeeper"
)

var secret string
var str string

var cmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypt a string",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	cmd.Flags().StringVarP(&secret, "token", "t", "TOKEN_SECRET", "the encrypted strings secret token")
	cmd.Flags().StringVarP(&str, "encrypted-string", "e", "", "the encrypted string")

	cryptkeepercli.Cmd.AddCommand(cmd)
}

func run() error {
	cryptkeeper.SetKey(secret)
	res, err := cryptkeeper.Encrypt(str)

	var style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF")).
		Background(lipgloss.Color("#158978")).
		Bold(true).
		MarginTop(1).
		PaddingLeft(4).
		PaddingRight(4)

	fmt.Println(style.Render(fmt.Sprintf(`
hash: %v
salt: %v
`, res.Hash, res.Salt)))

	return err
}
