package token

import (
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("TOKEN_SECRET", "wombat")
}

func TestTokenGenAndTest(t *testing.T) {
	tokenStr := Gen(Claims{"hi": "there"}, 0)
	token, err := Parse(tokenStr)

	assert.Equal(t, true, token.Valid)
	assert.Equal(t, "there", token.Claims.(jwt.MapClaims)["hi"].(string))

	err = Test(tokenStr)

	assert.NoError(t, err)
}
