package kafka

type Result struct {
	Key     string
	Content []interface{}
}

type Processor func([]byte) (*Result, error)
