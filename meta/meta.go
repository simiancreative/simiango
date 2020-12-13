package meta

import (
	"github.com/google/uuid"
)

type RequestId string

func Id() RequestId {
	uuid := uuid.New()
	return RequestId(uuid.String())
}
