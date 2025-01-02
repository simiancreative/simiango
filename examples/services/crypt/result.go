package crypt

import (
	"encoding/json"

	"github.com/simiancreative/simiango/cryptkeeper"
	"github.com/simiancreative/simiango/cryptkeeper/keepers/aes"
)

type OutgoingSettings struct {
	Email         string        `json:"email"`
	SensitiveData aes.Container `json:"sensitive"`
}

type IncomingSettings struct {
	Email         string        `json:"email"`
	SensitiveData aes.Container `json:"sensitive"`
}

type cryptService struct{}

func (s *cryptService) Result() (interface{}, error) {
	result := []interface{}{}

	data := []byte(`{ "email": "123@me.com", "sensitive": { "any": "super private info" } }`)

	i := IncomingSettings{}
	i.SensitiveData.SetKey("TOKEN_SECRET")
	o := OutgoingSettings{}
	o.SensitiveData.SetKey("TOKEN_SECRET")

	json.Unmarshal(data, &i)
	encryptedData, _ := json.Marshal(i)
	json.Unmarshal(encryptedData, &o)

	result = append(result, "original data")
	result = append(result, string(data))
	result = append(result, "encrypted using aes-256-cbc")
	result = append(result, i)
	result = append(result, "decrypted")
	result = append(result, o)

	keeper, err := cryptkeeper.New(cryptkeeper.BCRYPT)
	if err != nil {
		return nil, err
	}

	pass := "superSecret123123"
	enc, _ := keeper.Hash(pass)

	result = append(result, "password")
	result = append(result, pass)
	result = append(result, "encypted password using non reversible brcypt")
	result = append(result, enc)
	result = append(result, "compare password")
	result = append(result, keeper.Verify(pass, enc))

	return result, nil
}
