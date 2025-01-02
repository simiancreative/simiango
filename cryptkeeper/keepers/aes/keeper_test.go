package aes_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/simiancreative/simiango/cryptkeeper/keepers"
	"github.com/simiancreative/simiango/cryptkeeper/keepers/aes"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	os.Setenv("MY_SECRET_KEY", "my_secret_key_value")
	keeper := aes.Keeper{}
	config := keeper.Setup("MY_SECRET_KEY")
	assert.NotNil(t, config)
	assert.Equal(t, "my_secret_key_value", config.Get("key"))
}

func TestEncrypt(t *testing.T) {
	keeper := aes.Keeper{}
	// base64 encoded pass phrase which needs to be 32 bytes long
	config := keepers.Config{"key": "dGhpc2lzYXBhc3NrZXl3aGljaG5lZWRzMzJieXRlc3M="}

	content := strings.NewReader("my secret content")
	encrypted, err := keeper.Encrypt(config, content)

	assert.NoError(t, err)
	assert.NotNil(t, encrypted)
}

func TestDecrypt(t *testing.T) {
	keeper := aes.Keeper{}
	config := keepers.Config{}
	config.Set("key", "dGhpc2lzYXBhc3NrZXl3aGljaG5lZWRzMzJieXRlc3M=")

	content := strings.NewReader("my secret content")
	encrypted, err := keeper.Encrypt(config, content)
	assert.NoError(t, err)
	assert.NotNil(t, encrypted)

	decrypted, err := keeper.Decrypt(config, encrypted)
	assert.NoError(t, err)
	assert.NotNil(t, decrypted)

	decryptedContent, err := io.ReadAll(decrypted)
	assert.NoError(t, err)
	assert.Equal(t, "my secret content", string(decryptedContent))
}
