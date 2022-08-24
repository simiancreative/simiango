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

func (ps RawParams) AsMap() map[string]string {
	result := map[string]string{}

	for _, entry := range ps {
		result[entry.Key] = entry.Value
	}

	return result
}

func (ps RawParams) Assign(v interface{}) error {
	return parseParam("param", v, ps)
}
