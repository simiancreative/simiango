package main

import (
	_ "github.com/simiancreative/simiango/cmd/app/gen/service"

	_ "github.com/simiancreative/simiango/cmd/app/gen"

	"github.com/simiancreative/simiango/cmd/app"
)

func main() {
	app.Execute()
}
