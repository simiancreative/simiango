package sample

import (
	"encoding/json"

	"github.com/simiancreative/simiango/meta"
)

// SampleResource represents the model for an order
type sampleResource struct {
	Wibble string `json:"wibble" validate:"nonzero"`
}

type sampleResp struct {
	RequestID string `json:"request_id"`
	TokenID   string `json:"token_id"`
}

func (r *sampleResp) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}

type sampleContext struct {
	ID          string `json:"id"`
	Descendants string `json:"descendantsOf"`
}

type sampleHeaders map[string]string

type sampleService struct {
	id      meta.RequestId
	body    sampleResource
	params  sampleContext
	headers sampleHeaders
}
