package service

type RawHeaders []RawHeader
type RawHeader struct {
	Key    string
	Value  string
	Values []string
}

func (ps RawHeaders) Get(name string) (RawHeader, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			entry.Value = entry.Values[0]
			return entry, true
		}
	}
	return RawHeader{}, false
}
