package sample

import (
	"github.com/simiancreative/simiango/meta"
)

type sampleResource struct {
	Wibble string `json:"wibble"`
}

type sampleContext struct {
	ID         string `json:"id"`
	Decendants string `json:"decendantsOf"`
}

type sampleService struct {
	id     meta.RequestId
	body   sampleResource
	params sampleContext
}
