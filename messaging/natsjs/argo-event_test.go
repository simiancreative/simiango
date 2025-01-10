package natsjs_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/simiancreative/simiango/messaging/natsjs"
	"github.com/stretchr/testify/assert"
)

func TestArgoEventFromString(t *testing.T) {
	testEvent := natsjs.ArgoEvent{
		Context: natsjs.ArgoEventContext{
			ID:              "testID",
			Source:          "testSource",
			SpecVersion:     "testSpecVersion",
			Type:            "testType",
			DataContentType: "testDataContentType",
			Subject:         "testSubject",
			Time:            "testTime",
		},
		Data: "testData",
	}

	bytes, _ := json.Marshal(testEvent)
	str := string(bytes)

	event, err := natsjs.ArgoEventFromString(str)

	assert.NoError(t, err)
	assert.Equal(t, testEvent.Context.ID, event.Context.ID)

	// Test with invalid JSON
	_, err = natsjs.ArgoEventFromString("{invalid json}")
	assert.Error(t, err)
}

func TestJSONBody(t *testing.T) {
	testData := natsjs.ArgoEventData{
		Subject: "testSubject",
		Body:    "testBody",
	}

	bytes, _ := json.Marshal(testData)
	str := base64.StdEncoding.EncodeToString(bytes)

	testEvent := natsjs.ArgoEvent{
		Context: natsjs.ArgoEventContext{
			ID:              "testID",
			Source:          "testSource",
			SpecVersion:     "testSpecVersion",
			Type:            "testType",
			DataContentType: "testDataContentType",
			Subject:         "testSubject",
			Time:            "testTime",
		},
		Data: str,
	}

	body, err := testEvent.JSONBody()

	assert.NoError(t, err)
	assert.Equal(t, `"testBody"`, body)

	// Test with invalid base64 data
	testEvent.Data = "invalid base64"

	_, err = testEvent.JSONBody()
	assert.Error(t, err)
}

func TestBadJSONBody(t *testing.T) {
	str := base64.StdEncoding.EncodeToString([]byte("invalidJSON"))

	testEvent := natsjs.ArgoEvent{
		Context: natsjs.ArgoEventContext{
			ID:              "testID",
			Source:          "testSource",
			SpecVersion:     "testSpecVersion",
			Type:            "testType",
			DataContentType: "testDataContentType",
			Subject:         "testSubject",
			Time:            "testTime",
		},
		Data: str,
	}

	body, err := testEvent.JSONBody()
	assert.Error(t, err)
	assert.Equal(t, "", body)
}

func TestArgoUnmarshalEvent(t *testing.T) {
	// Sample JSON string representing an ArgoEvent
	sampleJSON := `{
		"context": {
			"id": "1234",
			"source": "source",
			"specversion": "1.0",
			"type": "type",
			"datacontenttype": "application/json",
			"subject": "subject",
			"time": "2023-10-01T00:00:00Z"
		},
		"data": "eyJzdWJqZWN0IjoiZXhhbXBsZSIsImJvZHkiOnsia2V5IjoidmFsdWUifX0="
	}`

	// Destination variable
	var dest map[string]interface{}

	// Call the function
	err := natsjs.ArgoUnmarshalEvent(&dest, sampleJSON)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the results
	expectedBody := map[string]interface{}{"key": "value"}
	if dest["key"] != expectedBody["key"] {
		t.Errorf("Expected body %v, got %v", expectedBody, dest)
	}
}

func TestArgoUnmarshalEvent_Errors(t *testing.T) {
	// Test case 1: Invalid JSON string
	invalidJSON := `{"context": { "id": "1234" "source": "source" }}` // Missing comma
	var dest1 map[string]interface{}
	err := natsjs.ArgoUnmarshalEvent(&dest1, invalidJSON)
	assert.Error(t, err)

	// Test case 2: Invalid base64 data
	invalidBase64 := `{
		"context": {
			"id": "1234",
			"source": "source",
			"specversion": "1.0",
			"type": "type",
			"datacontenttype": "application/json",
			"subject": "subject",
			"time": "2023-10-01T00:00:00Z"
		},
		"data": "invalid_base64_data"
	}`
	var dest2 map[string]interface{}
	err = natsjs.ArgoUnmarshalEvent(&dest2, invalidBase64)
	assert.Error(t, err)

	// test case 3: invalid json body
	invalidjsonbody := `{
		"context": {
			"id": "1234",
			"source": "source",
			"specversion": "1.0",
			"type": "type",
			"datacontenttype": "application/json",
			"subject": "subject",
			"time": "2023-10-01t00:00:00z"
		},
		"data": "W30K" // {"subject":"example","body":{"foo":"bar"}}
	}`
	var dest3 map[string]interface{}
	err = natsjs.ArgoUnmarshalEvent(&dest3, invalidjsonbody)
	assert.Error(t, err)

	// test case 3: invalid json body
	invalidjsonMessage := `{
		"context": {
			"id": "1234",
			"source": "source",
			"specversion": "1.0",
			"type": "type",
			"datacontenttype": "application/json",
			"subject": "subject",
			"time": "2023-10-01t00:00:00z"
		},
		"data": "eyJzdWJqZWN0IjoiZXhhbXBsZSIsImJvZHkiOiJ7XSJ9Cg=="
	}`
	var dest4 map[string]interface{}
	err = natsjs.ArgoUnmarshalEvent(&dest4, invalidjsonMessage)
	assert.Error(t, err)
}
