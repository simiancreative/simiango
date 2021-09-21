package crypt

import (
	"encoding/json"
	"github.com/simiancreative/simiango/cryptkeeper"
)

type OutgoingSettings struct {
	Email         string                           `json:"email"`
	SensitiveData cryptkeeper.DecryptableContainer `json:"sensitive"`
}

type IncomingSettings struct {
	Email         string                           `json:"email"`
	SensitiveData cryptkeeper.EncryptableContainer `json:"sensitive"`
}

type cryptService struct{}

func (s *cryptService) Result() (interface{}, error) {
	result := []interface{}{}

	data := []byte(`{ "email": "123@me.com", "sensitive": { "any": "super private info" } }`)
	i := IncomingSettings{}
	o := OutgoingSettings{}
	json.Unmarshal(data, &i)
	encryptedData, _ := json.Marshal(i)
	json.Unmarshal(encryptedData, &o)

	result = append(result, "original data")
	result = append(result, string(data))
	result = append(result, "encrypted using aes-256-cbc")
	result = append(result, i)
	result = append(result, "decrypted")
	result = append(result, o)

	pass := "superSecret123123"
	enc, _ := cryptkeeper.EncryptPassword(pass)

	result = append(result, "password")
	result = append(result, pass)
	result = append(result, "encypted password using non reversible brcypt")
	result = append(result, enc)
	result = append(result, "compare password")
	result = append(result, cryptkeeper.ComparePassword(pass, enc))

	return result, nil
}
