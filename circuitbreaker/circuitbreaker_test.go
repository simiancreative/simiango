package circuitbreaker_test

import (
	"os"
	"testing"
	"time"

	"github.com/simiancreative/simiango/circuitbreaker"
	"github.com/simiancreative/simiango/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type actionType int

const (
	allow actionType = iota
	start
	recordSuccess
	recordFailure
	checkState
	wait
)

type step struct {
	action   actionType
	want     interface{}   // bool for allow/start, State for checkState
	duration time.Duration // for wait action
}

func TestCircuitBreaker(t *testing.T) {
	os.Setenv("LOG_TYPE", "line")
	logger.Enable()

	tests := []struct {
		name           string
		config         circuitbreaker.Config
		steps          []step
		wantFinalState circuitbreaker.State
	}{
		{
			name: "remains closed when under threshold",
			config: circuitbreaker.Config{
				FailureThreshold: 3,
				OpenTimeout:      time.Second,
				HalfOpenMaxCalls: 2,
			},
			steps: []step{
				{action: start, want: true},
				{action: recordFailure},
				{action: start, want: true},
				{action: recordFailure},
				{action: allow, want: true},
				{action: checkState, want: circuitbreaker.StateClosed},
			},
			wantFinalState: circuitbreaker.StateClosed,
		},
		{
			name: "opens after threshold failures",
			config: circuitbreaker.Config{
				FailureThreshold: 2,
				OpenTimeout:      time.Second,
				HalfOpenMaxCalls: 2,
			},
			steps: []step{
				{action: start, want: true},
				{action: recordFailure},
				{action: start, want: true},
				{action: recordFailure},
				{action: start, want: false},
				{action: checkState, want: circuitbreaker.StateOpen},
			},
			wantFinalState: circuitbreaker.StateOpen,
		},
		{
			name: "transitions to half-open after timeout",
			config: circuitbreaker.Config{
				FailureThreshold: 2,
				OpenTimeout:      50 * time.Millisecond,
				HalfOpenMaxCalls: 2,
			},
			steps: []step{
				{action: start, want: true},
				{action: recordFailure},
				{action: start, want: true},
				{action: recordFailure},
				{action: wait, duration: 60 * time.Millisecond},
				{action: start, want: true},
				{action: checkState, want: circuitbreaker.StateHalfOpen},
			},
			wantFinalState: circuitbreaker.StateHalfOpen,
		},
		{
			name: "closes after successful half-open calls",
			config: circuitbreaker.Config{
				FailureThreshold: 2,
				OpenTimeout:      50 * time.Millisecond,
				HalfOpenMaxCalls: 2,
			},
			steps: []step{
				{action: start, want: true},
				{action: recordFailure},
				{action: start, want: true},
				{action: recordFailure},
				{action: wait, duration: 60 * time.Millisecond},
				{action: start, want: true},
				{action: recordSuccess},
				{action: start, want: true},
				{action: recordSuccess},
				{action: checkState, want: circuitbreaker.StateClosed},
			},
			wantFinalState: circuitbreaker.StateClosed,
		},
		{
			name: "limits calls in half-open state",
			config: circuitbreaker.Config{
				FailureThreshold: 2,
				OpenTimeout:      50 * time.Millisecond,
				HalfOpenMaxCalls: 2,
			},
			steps: []step{
				{action: start, want: true},
				{action: recordFailure},
				{action: start, want: true},
				{action: recordFailure},
				{action: wait, duration: 60 * time.Millisecond},
				{action: start, want: true},  // First call allowed
				{action: start, want: true},  // Second call allowed
				{action: start, want: false}, // Third call rejected
				{action: checkState, want: circuitbreaker.StateHalfOpen},
				{action: recordSuccess},
				{action: recordSuccess},
				{action: checkState, want: circuitbreaker.StateClosed},
			},
			wantFinalState: circuitbreaker.StateClosed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb, err := circuitbreaker.New(tt.config)
			require.NoError(t, err)

			for i, step := range tt.steps {
				switch step.action {
				case allow:
					assert.Equal(t, step.want, cb.Allow(), "step %d: Allow()", i)
				case start:
					assert.Equal(t, step.want, cb.RecordStart(), "step %d: RecordStart()", i)
				case recordSuccess:
					cb.RecordResult(true)
				case recordFailure:
					cb.RecordResult(false)
				case checkState:
					assert.Equal(t, step.want, cb.GetState(), "step %d: GetState()", i)
				case wait:
					time.Sleep(step.duration)
				}
			}

			assert.Equal(t, tt.wantFinalState, cb.GetState(), "final state")
		})
	}
}
