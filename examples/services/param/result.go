package param

import "github.com/sanity-io/litter"

func (s *paramService) Result() (interface{}, error) {
	litter.Dump(s)

	return s, nil
}
