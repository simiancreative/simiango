package assign

type assignService struct {
	event Event
}

func (s assignService) Result() (interface{}, error) {
	event := s.event
	assignable := AssignableStruct{}
	event.Body.Assign(&assignable)

	return assignable, nil
}
