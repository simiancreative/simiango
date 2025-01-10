package natsjs_test

import (
	"testing"

	"github.com/simiancreative/simiango/messaging/natsjs"
)

type TestEvent struct {
	Name string `json:"name"`
}

func TestUnmarshalEvent(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		wantEvent TestEvent
	}{
		{
			name:      "valid JSON",
			input:     `{"name": "test event"}`,
			wantErr:   false,
			wantEvent: TestEvent{Name: "test event"},
		},
		{
			name:    "invalid JSON",
			input:   `{"name": "test event"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var event TestEvent
			err := natsjs.UnmarshalEvent(&event, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && event != tt.wantEvent {
				t.Errorf("UnmarshalEvent() = %v, want %v", event, tt.wantEvent)
			}
		})
	}
}
