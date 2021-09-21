package cryptkeeper

import (
	"encoding/json"
	// "github.com/sanity-io/litter"
	"reflect"
)

type DecryptableData struct {
	Hash string `json:"hash"`
	Salt string `json:"salt"`
}

func (e DecryptableData) Decrypt() string {
	if e.Hash == "" && e.Salt == "" {
		return ""
	}

	result, _ := Decrypt(e.Hash, e.Salt)
	return result
}

type DecryptableContainer map[string]DecryptableData

func (dst *DecryptableContainer) Scan(src interface{}) error {
	if src == nil {
		*dst = DecryptableContainer{}
		return nil
	}

	v, _ := src.([]byte)
	t := DecryptableContainer{}
	json.Unmarshal(v, &t)

	*dst = t

	return nil
}

func (dst DecryptableContainer) MarshalJSON() ([]byte, error) {
	result := map[string]string{}

	keys := reflect.ValueOf(dst).MapKeys()
	for _, key := range keys {
		e := dst[key.String()]
		result[key.String()] = e.Decrypt()
	}

	return json.Marshal(result)
}
