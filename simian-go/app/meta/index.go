package metacli

import (
	"github.com/spf13/cobra"

	"github.com/simiancreative/simiango/simian-go/app"
)

var Cmd = &cobra.Command{
	Use: "meta",
}

func init() {
	app.RootCmd.AddCommand(Cmd)
}
