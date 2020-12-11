package sample

// import "errors"

func (s sampleService) Result() (interface{}, error) {
	return []interface{}{
		s.body,
		s.params,
	}, nil

	// example error response
	// return nil, errors.New("wooble")
}
