package decode

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	tokencli "github.com/simiancreative/simiango/simian-go/app/token"
	"github.com/simiancreative/simiango/token"
)

var key string
var tokenStr string

var cmd = &cobra.Command{
	Use:   "decode",
	Short: "decode a jwt token",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	cmd.Flags().StringVarP(&key, "key", "k", "", "the jwt secret key")
	cmd.Flags().StringVarP(&tokenStr, "token", "t", "", "the jwt token")

	tokencli.Cmd.AddCommand(cmd)
}

func run() error {
	if len(key) == 0 {
		return fmt.Errorf("secret required")
	}

	res, err := token.ParseWithSecret(
		tokenStr,
		[]byte(key),
	)

	if err != nil {
		return err
	}

	var style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF")).
		Background(lipgloss.Color("#158978")).
		Bold(true).
		MarginTop(1).
		PaddingTop(1).
		PaddingRight(4).
		PaddingBottom(1).
		PaddingLeft(4)

	resStr, err := json.MarshalIndent(res, "", "  ")
	fmt.Println(style.Render(string(resStr)))

	return nil
}
