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

//go:embed support/Reguard_mock_contracts.csv.pgp
var encStr string

//go:embed support/Reguard_mock_private.asc
var secretKeyring string

var passphrase = "reguard-mock-key"

func TestPGP(t *testing.T) {
	os.Setenv("PREFIX_PRIVATE_KEY", secretKeyring)
	os.Setenv("PREFIX_PASSPHRASE", passphrase)

	keeper := pgp.Keeper{}

	config := keeper.Setup("prefix")

	content, err := keeper.Decrypt(
		config,
		strings.NewReader(encStr),
	)

	assert.NoError(t, err)

	readContent, err := io.ReadAll(content)

	assert.NoError(t, err)
	assert.Equal(t, len(readContent), 1860)
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
