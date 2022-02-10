package service

type RawParams []ParamItem

func (ps RawParams) Get(name string) (ParamItem, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			return entry, true
		}
	}
	return ParamItem{}, false
}

func (ps RawParams) Assign(v interface{}) error {
	return parseParam("param", v, ps)
}
