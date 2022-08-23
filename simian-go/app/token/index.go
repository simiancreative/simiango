package tokencli

import (
	"github.com/spf13/cobra"

	"github.com/simiancreative/simiango/simian-go/app"
)

var Cmd = &cobra.Command{
	Use: "token",
}

func init() {
	app.RootCmd.AddCommand(Cmd)
}
