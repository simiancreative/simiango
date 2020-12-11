package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToResultError(t *testing.T) {
	err := ToResultError(errors.New("betty crocker"), "the cake is a lie", 500)

	assert.Equal(t, "betty crocker", err.Error())
	assert.Equal(t, "the cake is a lie", err.GetMessage())
	assert.Equal(t, 500, err.GetStatus())

	reasons := []map[string]interface{}{
		{"message": err.Error()},
	}

	assert.Equal(t, reasons, err.GetDetails()["reasons"])
}
