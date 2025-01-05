package pgp

import (
	"encoding/base64"
	"io"
	"os"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/simiancreative/simiango/cryptkeeper/keepers"
)

type Keeper struct{}

func (p Keeper) Setup(input ...any) keepers.Config {
	if len(input) == 0 {
		return nil
	}

	varPrefix, ok := input[0].(string)
	if !ok {
		return nil
	}

	config := keepers.Config{}

	// load keys and passphrase from env vars
	config.
		Set("privateKey", os.Getenv(strings.ToUpper(varPrefix)+"_PRIVATE_KEY")).
		Set("passphrase", os.Getenv(strings.ToUpper(varPrefix)+"_PASSPHRASE")).
		Set("publicKey", os.Getenv(strings.ToUpper(varPrefix)+"_PUBLIC_KEY"))

	return config
}

func (p Keeper) Decrypt(config keepers.Config, file io.Reader) (io.Reader, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return file, err
	}

	privateKey, _ := base64.StdEncoding.DecodeString(config.Get("privateKey"))

	decrypted, err := helper.DecryptMessageArmored(
		string(privateKey),
		[]byte(config.Get("passphrase")),
		string(content),
	)

	if err != nil {
		return file, err
	}

	return strings.NewReader(decrypted), nil
}

func (p Keeper) Encrypt(config keepers.Config, reader io.Reader) (io.Reader, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return reader, err
	}

	publicKey, _ := base64.StdEncoding.DecodeString(config.Get("publicKey"))
	privateKey, _ := base64.StdEncoding.DecodeString(config.Get("privateKey"))

	encrypted, err := helper.EncryptSignMessageArmored(
		string(publicKey),
		string(privateKey),
		[]byte(config.Get("passphrase")),
		string(content),
	)

	return strings.NewReader(encrypted), err
}

func New() *keepers.Keeper {
	p := Keeper{}
	k := &keepers.Keeper{}

	return k.
		SetConfigurator(p).
		SetEncrypter(p)
}
