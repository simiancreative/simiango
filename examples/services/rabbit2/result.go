package rabbit

import ()

type rabbitService struct {
	Hello string `json:"hello"`
}

func (s *rabbitService) Result() (interface{}, error) {
	return nil, nil
}
