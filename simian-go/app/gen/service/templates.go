package genservice

var auth = `package {{ServiceName}}

import (
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

func (s *{{ServiceName}}Service) Auth(
	requestID meta.RequestId,
	rawHeaders service.RawHeaders,
	rawBody service.RawBody,
	rawParams service.RawParams,
) bool {
	return true
}`

var authTest = `package {{ServiceName}}

import (
	"fmt"
	"testing"

	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
	"github.com/simiancreative/simiango/token"

	"github.com/stretchr/testify/assert"
)

var testTokenStr string

func init() {
	testTokenStr = token.GenWithSecret(token.Claims{
	}, bytes("key"), 0)
}

func TestAuth(t *testing.T) {
	var requestID meta.RequestId = meta.Id()
	var rawBody service.RawBody = []byte("")
	var rawParams = service.RawParams{}
	var rawHeaders = service.RawHeaders{
		{
			Key:    "Authorization",
			Values: []string{fmt.Sprintf("Bearer %s", testTokenStr)},
		},
	}

	service := &{{ServiceName}}Service{}
	ok := service.Auth(requestID, rawHeaders, rawBody, rawParams)

	assert.Equal(t, ok, true)
}`

var build = `package {{ServiceName}}

import (
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

func Build(
	requestID meta.RequestId,
	rawHeaders service.RawHeaders,
	rawBody service.RawBody,
	rawParams service.RawParams,
) (
	service.TPL, error,
) {
	// create service instance
	service := {{ServiceName}}Service{}
	return &service, nil
}`

var buildTest = `package {{ServiceName}}

import (
	"testing"

	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	requestID := meta.Id()
	rawBody := service.RawBody("")
	rawParams := service.RawParams{}
	rawHeaders := service.RawHeaders{}

	service, err := Build(requestID, rawHeaders, rawBody, rawParams)

	assert.NoError(t, err)
	assert.IsType(t, service, &{{ServiceName}}Service{})
}`

var config = `package {{ServiceName}}

import (
	"github.com/simiancreative/simiango/service"
	"github.com/simiancreative/simiango/server"
)

{{#if SetupSwagger}}
// Build godoc
// @Summary [enter service summary]
// @Description [enter service description]
// @Tags [add tags]
// @ID {{ServiceName}}
// @Produce  json
// @Success [replace with valid response code] {object} [replace with response struct]
// @Failure [replace with error code] {object} service.ResultError "[replace with error_name]"
{{#if IsPrivate}}
// @Security bearer_auth
{{/if}}
// @Router /{{ServiceURL}} [{{ServiceMethod}}]
{{/if}}
var Config = service.Config{
	IsPrivate: {{IsPrivate}},
	Method:    "{{ServiceMethod}}",
	Path:      "/{{ServiceURL}}",
	Build:     Build,
}

// dont forget to import your package in your main.go for initialization 
// _ "path/to/project/{{ServiceName}}"
func init() {
	server.AddService(Config)
}`

var configTest = `package {{ServiceName}}

import (
	"testing"

	"github.com/simiancreative/simiango/service"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert.Equal(t, Config.Path, "/{{ServiceURL}}")
	assert.Equal(t, Config.Method, "{{ServiceMethod}}")
	assert.IsType(t, Config.Build, service.Config{}.Build)
}`

var result = `package {{ServiceName}}

type {{ServiceName}}Service struct {}

func (s *{{ServiceName}}Service) Result() (interface{}, error) {
	return "", nil
}`

var resultTest = `package {{ServiceName}}

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	s := &{{ServiceName}}Service{}
	result, err := s.Result()

	assert.Equal(t, err, nil)
	assert.Equal(t, result, "")
}`
