package assign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testassign struct {
	Body Assigner
}

type testbody struct {
	Hi string `json:"hi" assign:"hi"`
}

func TestAssignable(t *testing.T) {
	obj := testassign{Body: map[string]interface{}{
		"hi": "there",
	}}

	receiver := testbody{}
	obj.Body.Assign(&receiver)

	assert.Equal(t, receiver.Hi, "there")
}
