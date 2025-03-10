package metacli

import (
	"github.com/spf13/cobra"

	"github.com/simiancreative/simiango/simian-go/app"
)

var Cmd = &cobra.Command{
	Use: "tryitout",
}

func init() {
	app.Root.Cmd.AddCommand(Cmd)
}
