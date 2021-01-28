package configdata

import (
	"github.com/simiancreative/simiango/meta"
)

type configDataService struct {
	id	meta.RequestId
}

type ConfigDataResp struct {
	Tag string `json:"tag" db:"Tag"`
	PlatformType string `json:"platformtype" db:"PlatformType"`
	AppHash string `json:"apphash" db:"AppHash"`
	TokenSigningKey string `json:"tokensigningkey" db:"TokenSigningKey"`
	RegistrationExpMilliseconds int `json:"registrationexpmilliseconds" db:"RegistrationExpMilliseconds"`
	TokenExpMinutesSeconds int `json:"tokenexpminutesseconds" db:"TokenExpMinutesSeconds"`
	StatusID int `json:"statusid" db:"StatusID"`
}