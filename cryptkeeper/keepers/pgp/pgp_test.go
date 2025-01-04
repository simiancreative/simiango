package pgp_test

import (
	_ "embed"

	"io"
	"os"
	"strings"
	"testing"

	"github.com/simiancreative/simiango/cryptkeeper/keepers/pgp"
	"github.com/stretchr/testify/assert"
)

//go:embed support/Reguard_mock_contracts.csv
var str string

//go:embed support/Reguard_mock_private.asc
var secretKeyring string

//go:embed support/Reguard_mock_public.asc
var publicKeyring string

var passphrase = "reguard-mock-key"

func TestPGP(t *testing.T) {
	os.Setenv("PREFIX_PUBLIC_KEY", publicKeyring)
	os.Setenv("PREFIX_PRIVATE_KEY", secretKeyring)
	os.Setenv("PREFIX_PASSPHRASE", passphrase)

	keeper := pgp.New()

	keeper.Setup("prefix")

	encrypted, err := keeper.Encrypt(
		strings.NewReader(str),
	)

	assert.NotEqual(t, encrypted, strings.NewReader(str))
	assert.NoError(t, err)

	content, err := keeper.Decrypt(
		encrypted,
	)

	assert.NoError(t, err)

	readContent, err := io.ReadAll(content)

	assert.NoError(t, err)
	assert.Equal(t, len(readContent), 1860)
	assert.Equal(t, string(readContent), str)
}

func TestPGPErr(t *testing.T) {
	os.Setenv("PREFIX_PRIVATE_KEY", secretKeyring)
	os.Setenv("PREFIX_PASSPHRASE", passphrase)

	keeper := pgp.Keeper{}

	config := keeper.Setup("prefix")

	_, err := keeper.Decrypt(
		config,
		strings.NewReader("bad"),
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "EOF")
}
