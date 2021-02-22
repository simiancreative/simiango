package service

type RawHeaders []RawHeader
type RawHeader struct {
	Key    string
	Values []string
}

func (ps RawHeader) Value() string {
	return ps.Values[0]
}

func (ps RawHeaders) Get(name string) (*RawHeader, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			return &entry, true
		}
	}
	return nil, false
}
