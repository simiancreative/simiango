package nats_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/simiancreative/simiango/messaging/natsjs"
	"github.com/simiancreative/simiango/mocks/nats"
	"github.com/stretchr/testify/assert"
)

func TestBuildArgoEvent(t *testing.T) {
	// Test the happy path
	workflowName := "test-workflow"
	action := "test-action"
	msg := `{"key":"value"}`

	result := nats.BuildArgoEvent(workflowName, action, msg)

	var event natsjs.ArgoEvent
	err := json.Unmarshal([]byte(result), &event)
	assert.NoError(t, err)

	expectedSubject := workflowName + "-" + action
	assert.Equal(t, expectedSubject, event.Context.Subject)

	var data natsjs.ArgoEventData
	dataBytes, _ := base64.StdEncoding.DecodeString(event.Data)
	json.Unmarshal(dataBytes, &data)
	assert.Equal(t, data.Subject, expectedSubject)
}
