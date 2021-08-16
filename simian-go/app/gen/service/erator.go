package genservice

import (
	"github.com/simiancreative/simiango/cmd/app/gen"
)

var erator = gen.Erator{
	Name: "service",
	Desc: "package defining a simian-go service",
	RequiredVars: []gen.RequiredVar{
		{
			Type:    "string",
			Name:    "ServiceName",
			Message: "Enter the services name",
		},
		{
			Type:    "string",
			Name:    "ServiceMethod",
			Message: "enter the service's http method",
		},
		{
			Type:    "string",
			Name:    "ServiceURL",
			Message: "Enter the services url including path parameters (/wibble/:wibbles)",
		},
		{
			Type:    "confirm",
			Name:    "IsPrivate",
			Message: "Does the service require Auth",
		},
	},
	Templates: []gen.Template{
		{
			IF:      "IsPrivate",
			Path:    "auth.go",
			Content: auth,
		},
		{
			IF:      "IsPrivate",
			Path:    "auth_test.go",
			Content: authTest,
		},
		{
			Path:    "build.go",
			Content: build,
		},
		{
			Path:    "build_test.go",
			Content: buildTest,
		},
		{
			Path:    "config.go",
			Content: config,
		},
		{
			Path:    "config_test.go",
			Content: configTest,
		},
		{
			Path:    "result.go",
			Content: result,
		},
		{
			Path:    "result_test.go",
			Content: resultTest,
		},
	},
}

func init() {
	gen.Register(erator)
}
