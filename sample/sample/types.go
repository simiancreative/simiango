package sample

import (
	"github.com/simiancreative/simiango/context"
)

type SampleResource struct {
	Wibble string `json:"wibble"`
}

type SampleContext struct {
	ID         string `json:"id"`
	Decendants string `json:"decendantsOf"`
}

type sampleService struct {
	id     context.RequestId
	body   SampleResource
	params SampleContext
}
