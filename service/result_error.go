package service

type ResultError struct {
	Status     int
	ErrMessage string
	Message    string
	Field      string
	Reasons    []map[string]interface{} `json:"details"`
}

func (m *ResultError) Error() string {
	return m.ErrMessage
}
func (m *ResultError) GetMessage() string {
	return m.Message
}
func (m *ResultError) GetStatus() int {
	return m.Status
}
func (m *ResultError) GetDetails() map[string]interface{} {
	return map[string]interface{}{
		"message": m.Message,
		"field":   m.Field,
		"reasons": m.Reasons,
	}
}

func ToResultError(err error, message string, status int) *ResultError {
	reasons := []map[string]interface{}{
		{"message": err.Error()},
	}

	return &ResultError{
		ErrMessage: err.Error(),
		Status:     status,
		Message:    message,
		Reasons:    reasons,
	}
}
