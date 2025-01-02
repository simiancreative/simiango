package decrypt

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/simiancreative/simiango/cryptkeeper"
	"github.com/simiancreative/simiango/cryptkeeper/keepers/aes"
	cryptkeepercli "github.com/simiancreative/simiango/simian-go/app/cryptkeeper"
)

var secret string
var hash string
var salt string

var cmd = &cobra.Command{
	Use:   "decrypt",
	Short: "decrypt a hash and salt",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	cmd.Flags().
		StringVarP(&secret, "token", "t", "TOKEN_SECRET", "the encrypted strings secret token")
	cmd.Flags().StringVarP(&hash, "hash", "a", "", "the hash")
	cmd.Flags().StringVarP(&salt, "salt", "s", "", "the salt")

	cryptkeepercli.Cmd.AddCommand(cmd)
}

func run() error {
	keeper, err := cryptkeeper.New(cryptkeeper.AES)
	if err != nil {
		return err
	}

	data := aes.Data{Hash: hash, Salt: salt}
	os.Setenv("AES_TOKEN", secret)

	res, err := keeper.Setup("AES_TOKEN").Decrypt(data)

	var style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF")).
		Background(lipgloss.Color("#158978")).
		Bold(true).
		MarginTop(1).
		PaddingLeft(4).
		PaddingRight(4)

	fmt.Println(style.Render(fmt.Sprintf(`
decrypted: %v
`, res)))

	return err
}
