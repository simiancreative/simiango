package service

type RawHeaders []ParamItem

func (ps RawHeaders) Get(name string) (ParamItem, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			entry.Value = entry.Values[0]
			return entry, true
		}
	}
	return ParamItem{}, false
}

func (ps RawHeaders) Assign(v interface{}) error {
	return parseParam("header", v, ps)
}
