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

func (ps RawParams) GetWithFallback(name string, fallback string) ParamItem {
	for _, entry := range ps {
		if entry.Key == name {
			return entry
		}
	}
	return ParamItem{Value: fallback}
}

func (ps RawParams) Assign(v interface{}) error {
	return parseParam("param", v, ps)
}
