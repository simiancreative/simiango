package aes

import (
	"database/sql/driver"
	"encoding/json"
	"os"
	"reflect"
)

type Container struct {
	Data map[string]Data
	Key  string
}

func (c *Container) SetKey(keyName string) *Container {
	c.Key = os.Getenv(keyName)

	return c
}

func (c *Container) Scan(src interface{}) error {
	if src == nil {
		c.Data = make(map[string]Data)
		return nil
	}

	v, _ := src.([]byte)
	t := Container{}
	err := json.Unmarshal(v, &t.Data)
	if err != nil {
		return err
	}

	*c = t

	return nil
}

func (c Container) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Container) UnmarshalJSON(b []byte) error {
	dst := Container{}
	src := map[string]string{}
	err := json.Unmarshal(b, &src)
	if err != nil {
		return err
	}

	keys := reflect.ValueOf(src).MapKeys()
	for _, key := range keys {
		e := src[key.String()]
		val, _ := encrypt(c.Key, e)
		c.Data[key.String()] = val
	}

	*c = dst

	return nil
}

func (c Container) MarshalJSON() ([]byte, error) {
	result := map[string]string{}

	keys := reflect.ValueOf(c).MapKeys()
	for _, key := range keys {
		e := c.Data[key.String()]
		value, err := decrypt(c.Key, e.Hash, e.Salt)
		if err != nil {
			return nil, err
		}

		result[key.String()] = value
	}

	return json.Marshal(result)
}
