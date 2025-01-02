package aes

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/simiancreative/simiango/cryptkeeper/keepers"
)

type Keeper struct{}

func (p Keeper) Setup(input ...any) keepers.Config {
	if len(input) == 0 {
		return nil
	}

	keyName, ok := input[0].(string)
	if !ok {
		return nil
	}

	config := keepers.Config{}

	config.Set("key", os.Getenv(keyName))

	return config
}

func (p Keeper) Encrypt(config keepers.Config, content io.Reader) (io.Reader, error) {
	value, err := io.ReadAll(content)
	if err != nil {
		return nil, err
	}

	return encrypt(config.Get("key"), string(value))
}

func (p Keeper) Decrypt(config keepers.Config, content io.Reader) (io.Reader, error) {
	value, err := io.ReadAll(content)
	if err != nil {
		return nil, err
	}

	dest := Data{}
	err = json.Unmarshal(value, &dest)
	if err != nil {
		return nil, err
	}

	result, err := decrypt(config.Get("key"), dest.Hash, dest.Salt)
	if err != nil {
		return nil, err
	}

	reader := strings.NewReader(result)

	return reader, nil
}

func New() *keepers.Keeper {
	p := Keeper{}
	k := &keepers.Keeper{}

	return k.
		SetConfigurator(p).
		SetEncrypter(p)
}
