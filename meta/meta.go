package meta

import (
	"time"

	"github.com/google/uuid"
)

type RequestId string

func Id() RequestId {
	uuid := uuid.New()
	return RequestId(uuid.String())
}

func GetDurationInMillseconds(start time.Time) float64 {
	end := time.Now()
	duration := end.Sub(start)
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	return rounded
}
