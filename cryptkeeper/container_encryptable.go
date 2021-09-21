package cryptkeeper

import (
	"database/sql/driver"
	"encoding/json"
	// "github.com/sanity-io/litter"
	"reflect"
)

type EncryptableData struct {
	Hash string `json:"hash"`
	Salt string `json:"salt"`
}

func (e EncryptableData) Decrypt() string {
	if e.Hash == "" && e.Salt == "" {
		return ""
	}

	result, _ := Decrypt(e.Hash, e.Salt)
	return result
}

type EncryptableContainer map[string]EncryptableData

func (dst EncryptableContainer) Value() (driver.Value, error) {
	return json.Marshal(dst)
}

func (dst *EncryptableContainer) UnmarshalJSON(b []byte) error {
	c := EncryptableContainer{}
	src := map[string]string{}
	json.Unmarshal(b, &src)

	keys := reflect.ValueOf(src).MapKeys()
	for _, key := range keys {
		e := src[key.String()]
		val, _ := Encrypt(e)
		c[key.String()] = val
	}

	*dst = c

	return nil
}
