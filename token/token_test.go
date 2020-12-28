package token

import (
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/simplereach/timeutils"
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

type tokenCarrier struct {
	JTI     string         `json:"jti"`
	IAT     timeutils.Time `json:"iat"`
	ANI     string         `json:"ani"`
	Brand   string         `json:"brand"`
	Version string         `json:"version"`
}

func TestTokenDecode(t *testing.T) {
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhbmkiOiIzMDU0MDkyMjI1IiwiYnJhbmQiOiJ0cmFjZm9uZSIsImlhdCI6MTYwODgzNjg0NSwianRpIjoiMTA0ODgzZDktYzY0Yy00YWEyLWFmYzMtODk5MDdhNTMyNzhhIiwidmVyc2lvbiI6ImVuaWdtYXRpYy13b21iYXQifQ.VOrbGkOgVEbgUpRS1gnNGQSABkdJw_wKx4vAGQC8m0w"
	carrier := &tokenCarrier{}
	err := Decode(tokenStr, carrier)

	assert.Equal(t, "tracfone", carrier.Brand)
	assert.Equal(t, "enigmatic-wombat", carrier.Version)
	assert.Equal(t, int64(1608836845), carrier.IAT.Unix())
	assert.NoError(t, err)
}
