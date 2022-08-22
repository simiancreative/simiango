package main

import (
	_ "github.com/simiancreative/simiango/simian-go/app/gen/service"

	_ "github.com/simiancreative/simiango/simian-go/app/cryptkeeper"
	_ "github.com/simiancreative/simiango/simian-go/app/cryptkeeper/decrypt"
	_ "github.com/simiancreative/simiango/simian-go/app/cryptkeeper/encrypt"
	_ "github.com/simiancreative/simiango/simian-go/app/gen"

	"github.com/simiancreative/simiango/simian-go/app"
)

func main() {
	app.Execute()
}
