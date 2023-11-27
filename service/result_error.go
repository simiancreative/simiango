package service

type ResultError struct {
	Status     int                      `json:"-"`
	Field      string                   `json:"-"`
	ErrMessage string                   `json:"-"`
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

type Reason map[string]interface{}

func Error(
	status int,
	err error,
	reasons ...Reason,
) *ResultError {
	maps := []map[string]interface{}{}

	for _, reason := range reasons {
		maps = append(maps, reason)
	}

	return &ResultError{
		Status:     status,
		ErrMessage: err.Error(),
		Message:    err.Error(),
		Reasons:    maps,
	}
}
