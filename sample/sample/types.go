package sample

import (
	"github.com/simiancreative/simiango/meta"
)

type SampleResource struct {
	Wibble string `json:"wibble"`
}

type SampleContext struct {
	ID         string `json:"id"`
	Decendants string `json:"decendantsOf"`
}

type sampleService struct {
	id     meta.RequestId
	body   SampleResource
	params SampleContext
}
