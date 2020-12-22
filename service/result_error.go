package service

type ResultError struct {
	Status     int
	Field      string
	ErrMessage string
	Message    string                   `json:"message"`
	Reasons    []map[string]interface{} `json:"reasons"`
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
func (m *ResultError) GetDetails() ResultError {
	return *m
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
