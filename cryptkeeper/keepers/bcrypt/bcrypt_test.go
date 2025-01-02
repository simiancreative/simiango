package bcrypt_test

import (
	"testing"

	"github.com/simiancreative/simiango/cryptkeeper/keepers"
	"github.com/simiancreative/simiango/cryptkeeper/keepers/bcrypt"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	keeper := bcrypt.Keeper{}
	password := "mysecretpassword"
	hash, err := keeper.Hash(keepers.Config{}, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestVerify(t *testing.T) {
	keeper := bcrypt.Keeper{}
	password := "mysecretpassword"
	hash, err := keeper.Hash(keepers.Config{}, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	valid := keeper.Verify(keepers.Config{}, hash, password)
	assert.True(t, valid)

	invalid := keeper.Verify(keepers.Config{}, hash, "wrongpassword")
	assert.False(t, invalid)
}
