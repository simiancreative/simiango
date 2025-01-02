package errors_test

import (
	"testing"

	"github.com/simiancreative/simiango/errors"
	"github.com/stretchr/testify/assert"
)

func TestReasons_Result(t *testing.T) {
	// Define some test items
	reasons := errors.Reasons{
		1: {Status: 200, Key: "OK", Description: "Operation successful"},
		2: {Status: 404, Key: "NotFound", Description: "Item not found"},
	}

	// Test case: Valid ID with no additional reasons
	result := reasons.Result(1)
	assert.NotNil(t, result)
	assert.Equal(t, 200, result.Status)
	assert.Equal(t, "Operation successful", result.ErrMessage)
	assert.Equal(t, "OK", result.Message)
	assert.Empty(t, result.Reasons)

	// Test case: Valid ID with additional reasons
	additionalReasons := errors.Reason{"detail": "extra info"}
	result = reasons.Result(2, additionalReasons)
	assert.NotNil(t, result)
	assert.Equal(t, 404, result.Status)
	assert.Equal(t, "Item not found", result.ErrMessage)
	assert.Equal(t, "NotFound", result.Message)
	assert.Len(t, result.Reasons, 1)
	assert.Equal(t, "extra info", result.Reasons[0]["detail"])

	// Test case: Invalid ID
	result = reasons.Result(999)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Status)
	assert.Empty(t, result.ErrMessage)
	assert.Empty(t, result.Message)
	assert.Empty(t, result.Reasons)
}
