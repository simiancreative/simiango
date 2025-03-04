package circuitbreaker

import (
	"fmt"
	"sync"
	"time"

	"github.com/simiancreative/simiango/logger"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

type Logger interface {
	Debug(...interface{})
}

type Config struct {
	FailureThreshold int
	OpenTimeout      time.Duration
	HalfOpenMaxCalls int
	OnStateChange    func(from, to State)
	Logger           Logger
}

type Breaker interface {
	Allow() bool
	GetState() State
	RecordStart() bool
	RecordResult(success bool)
	Reset()
}

type CircuitBreaker struct {
	config    Config
	state     State
	failures  int
	attempts  int
	successes int
	mutex     sync.RWMutex
	timer     *time.Timer
}

func NewDefault() (*CircuitBreaker, error) {
	return New(Config{
		FailureThreshold: 5,
		OpenTimeout:      5 * time.Second,
		HalfOpenMaxCalls: 5,
	})
}

func New(config Config) (*CircuitBreaker, error) {
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	breaker := &CircuitBreaker{
		config: config,
		state:  StateClosed,
	}

	breaker.config.Logger.Debug("creating new circuit breaker", logger.Fields{
		"failure_threshold":   config.FailureThreshold,
		"open_timeout":        config.OpenTimeout.String(),
		"half_open_max_calls": config.HalfOpenMaxCalls,
	})

	return breaker, nil
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	allowed := false
	switch cb.state {
	case StateOpen:
		allowed = false
	case StateHalfOpen:
		allowed = cb.attempts < cb.config.HalfOpenMaxCalls
	default:
		allowed = true
	}

	logger.Debug("circuit breaker allow check", logger.Fields{
		"state":     cb.state.String(),
		"allowed":   allowed,
		"attempts":  cb.attempts,
		"max_calls": cb.config.HalfOpenMaxCalls,
	})

	return allowed
}

func (cb *CircuitBreaker) GetState() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// RecordStart marks the beginning of an attempt
func (cb *CircuitBreaker) RecordStart() bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case StateOpen:
		logger.Debug("attempt rejected - circuit open", logger.Fields{
			"state": cb.state.String(),
		})
		return false
	case StateHalfOpen:
		if cb.attempts >= cb.config.HalfOpenMaxCalls {
			logger.Debug("attempt rejected - max half-open calls reached", logger.Fields{
				"attempts":  cb.attempts,
				"max_calls": cb.config.HalfOpenMaxCalls,
			})
			return false
		}
	}

	cb.attempts++
	cb.config.Logger.Debug("attempt started", logger.Fields{
		"state":    cb.state.String(),
		"attempts": cb.attempts,
	})
	return true
}

// RecordResult records the result of an attempt
func (cb *CircuitBreaker) RecordResult(success bool) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.config.Logger.Debug("recording attempt result", logger.Fields{
		"success":   success,
		"state":     cb.state.String(),
		"attempts":  cb.attempts,
		"successes": cb.successes,
		"failures":  cb.failures,
	})

	if !success {
		cb.recordFailure()
		return
	}

	switch cb.state {
	case StateHalfOpen:
		cb.successes++
		cb.config.Logger.Debug("recorded success in half-open state", logger.Fields{
			"attempts":  cb.attempts,
			"successes": cb.successes,
			"max_calls": cb.config.HalfOpenMaxCalls,
		})
		if cb.successes >= cb.config.HalfOpenMaxCalls {
			cb.transitionTo(StateClosed)
		}
	case StateClosed:
		cb.failures = 0
		cb.config.Logger.Debug("recorded success in closed state", logger.Fields{
			"failures": cb.failures,
		})
	}
}

func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.config.Logger.Debug("resetting circuit breaker", logger.Fields{
		"from_state": cb.state.String(),
	})

	if cb.timer != nil {
		cb.timer.Stop()
	}

	cb.transitionTo(StateClosed)
	cb.failures = 0
	cb.attempts = 0
	cb.successes = 0
}

func (cb *CircuitBreaker) recordFailure() {
	cb.failures++

	cb.config.Logger.Debug("recorded failure", logger.Fields{
		"state":     cb.state.String(),
		"failures":  cb.failures,
		"threshold": cb.config.FailureThreshold,
	})

	shouldOpen := cb.state == StateHalfOpen ||
		(cb.state == StateClosed && cb.failures >= cb.config.FailureThreshold)

	if !shouldOpen {
		return
	}

	cb.openCircuit()
}

func (cb *CircuitBreaker) openCircuit() {
	cb.config.Logger.Debug("opening circuit", logger.Fields{
		"from_state":   cb.state.String(),
		"open_timeout": cb.config.OpenTimeout.String(),
	})

	if cb.timer != nil {
		cb.timer.Stop()
	}

	cb.transitionTo(StateOpen)

	cb.timer = time.AfterFunc(cb.config.OpenTimeout, func() {
		cb.mutex.Lock()
		defer cb.mutex.Unlock()

		cb.config.Logger.Debug("open timeout elapsed", logger.Fields{
			"current_state": cb.state.String(),
		})

		if cb.state == StateOpen {
			cb.transitionTo(StateHalfOpen)
		}
	})
}

func (cb *CircuitBreaker) transitionTo(newState State) {
	if cb.state == newState {
		return
	}

	oldState := cb.state
	cb.state = newState
	cb.attempts = 0
	cb.successes = 0

	cb.config.Logger.Debug("state transition", logger.Fields{
		"from_state": oldState.String(),
		"to_state":   newState.String(),
		"attempts":   cb.attempts,
		"successes":  cb.successes,
	})

	if cb.config.OnStateChange != nil {
		go cb.config.OnStateChange(oldState, newState)
	}
}

func validateConfig(config *Config) error {
	if config.FailureThreshold <= 0 {
		return fmt.Errorf("failure threshold must be greater than 0")
	}

	if config.OpenTimeout <= 0 {
		return fmt.Errorf("open timeout must be greater than 0")
	}

	if config.HalfOpenMaxCalls <= 0 {
		return fmt.Errorf("half-open max calls must be greater than 0")
	}

	if config.Logger == nil {
		config.Logger = logger.New()
	}

	return nil
}
