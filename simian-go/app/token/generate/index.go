package generate

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	tokencli "github.com/simiancreative/simiango/simian-go/app/token"
	"github.com/simiancreative/simiango/token"
)

var key string
var exp int64
var claims string
var authorities string

var cmd = &cobra.Command{
	Use:   "generate",
	Short: "generate a jwt token",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	cmd.Flags().StringVarP(&key, "key", "k", "", "the jwt secret key")
	cmd.Flags().Int64VarP(&exp, "exp-minutes", "e", 0, "expiration in minutes from now")
	cmd.Flags().StringVarP(&claims, "claims", "c", "", "list of claims")
	cmd.Flags().StringVarP(&authorities, "authorities", "a", "", "list of authorities")

	tokencli.Cmd.AddCommand(cmd)
}

func run() error {
	if len(key) == 0 {
		return fmt.Errorf("secret required")
	}

	explodedClaims := token.Claims{}
	if len(claims) > 0 {
		params, err := url.ParseQuery(claims)
		if err != nil {
			return err
		}
		for key, value := range params {
			explodedClaims[key] = value[0]
		}
	}

	if len(authorities) > 0 {
		authorityStrings := strings.Split(authorities, ",")
		explodedClaims["authorities"] = authorityStrings
	}

	res := token.GenWithSecret(
		explodedClaims,
		[]byte(key),
		time.Duration(exp),
	)

	var style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF")).
		Background(lipgloss.Color("#158978")).
		Bold(true).
		MarginTop(1).
		PaddingTop(1).
		PaddingRight(4).
		PaddingBottom(1).
		PaddingLeft(4)

	fmt.Println(style.Render(fmt.Sprintf(`%v`, res)))

	return nil
}
