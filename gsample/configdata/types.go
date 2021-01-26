package configdata

import (
	"encoding/json"

	"github.com/simiancreative/simiango/meta"
	"github.com/simplereach/timeutils"
)

type configDataService struct {
	id	meta.RequestId
}

type ConfigDataResp struct {
	Tag string `json:"tag"`
	PlatformType string `json:"platformtype"`
	AppHash string `json:"apphash"`
	TokenSigningKey string `json:"tokensigningkey"`
	RegistrationExpMilliseconds int64 `json:"registrationexpmilliseconds"`
	TokenExpMinutesseconds int `json:"tokenexpminutesseconds"`
	StatusID int `json:"statusid"`
}