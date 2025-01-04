package workflow

import (
	"encoding/json"

	"github.com/simiancreative/simiango/errors"
)

type (
	ArgsList []ArgsItem
	ArgsItem []string
	Args     map[string]string
)

func (a Args) ItemFromJSON(key string) (result string) {
	result, ok := a[key]
	if !ok {
		return
	}

	json.Unmarshal([]byte(result), &result)

	return
}

func (a Args) UnmarshalKey(key string, dest interface{}) error {
	val, ok := a[key]
	if !ok {
		return errors.New("key not found in args %v", key)
	}

	return json.Unmarshal([]byte(val), dest)
}
