package service

type RawParams []RawParam
type RawParam struct {
	Key    string
	Value  string
	Values []string
}

func (ps RawParams) Get(name string) (*RawParam, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			return &entry, true
		}
	}
	return nil, false
}
