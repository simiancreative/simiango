package meta

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	res := Id()

	assert.IsType(t, RequestId(""), res)
}

func TestGetDurationMilliseconds(t *testing.T) {
	start := time.Now().Add(-time.Second * 10)

	dur := GetDurationInMillseconds(start)

	assert.IsType(t, dur, float64(0))
	assert.Greater(t, dur, float64(9999))
}
