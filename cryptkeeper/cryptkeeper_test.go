package cryptkeeper_test

import (
	"testing"

	"github.com/simiancreative/simiango/cryptkeeper"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	keeper, err := cryptkeeper.New(cryptkeeper.AES)
	assert.NoError(t, err)
	assert.NotNil(t, keeper)

	_, err = cryptkeeper.New(999) // Invalid kind
	assert.Error(t, err)
}
