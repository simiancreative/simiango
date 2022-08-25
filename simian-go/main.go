package main

import (
	_ "github.com/simiancreative/simiango/simian-go/app/gen/service"

	_ "github.com/simiancreative/simiango/simian-go/app/cryptkeeper"
	_ "github.com/simiancreative/simiango/simian-go/app/cryptkeeper/decrypt"
	_ "github.com/simiancreative/simiango/simian-go/app/cryptkeeper/encrypt"
	_ "github.com/simiancreative/simiango/simian-go/app/gen"
	_ "github.com/simiancreative/simiango/simian-go/app/meta"
	_ "github.com/simiancreative/simiango/simian-go/app/meta/uuid"
	_ "github.com/simiancreative/simiango/simian-go/app/token"
	_ "github.com/simiancreative/simiango/simian-go/app/token/decode"
	_ "github.com/simiancreative/simiango/simian-go/app/token/generate"
	_ "github.com/simiancreative/simiango/simian-go/app/token/test"

	"github.com/simiancreative/simiango/simian-go/app"
)

func main() {
	app.Execute()
}
