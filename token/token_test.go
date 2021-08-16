package token

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"

	"github.com/simiancreative/simiango/timeutils"
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

	tokenStr = Gen(Claims{"hi": "there"}, 300)
	token, _ = Parse(tokenStr)

	assert.Equal(t, true, token.Valid)
}

type tokenCarrier struct {
	JTI     string             `json:"jti"`
	IAT     timeutils.UnixTime `json:"iat"`
	ANI     string             `json:"ani"`
	Brand   string             `json:"brand"`
	Version string             `json:"version"`
}

func TestTokenDecode(t *testing.T) {
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhbmkiOiIzMDU0MDkyMjI1IiwiYnJhbmQiOiJ0cmFjZm9uZSIsImlhdCI6MTYwODgzNjg0NSwianRpIjoiMTA0ODgzZDktYzY0Yy00YWEyLWFmYzMtODk5MDdhNTMyNzhhIiwidmVyc2lvbiI6ImVuaWdtYXRpYy13b21iYXQifQ.VOrbGkOgVEbgUpRS1gnNGQSABkdJw_wKx4vAGQC8m0w"
	carrier := &tokenCarrier{}
	err := Decode(tokenStr, carrier)

	assert.Equal(t, "tracfone", carrier.Brand)
	assert.Equal(t, "enigmatic-wombat", carrier.Version)
	assert.Equal(t, int64(1608836845), carrier.IAT.Unix())
	assert.NoError(t, err)

	tokenStr = "1231234"
	carrier = &tokenCarrier{}
	err = Decode(tokenStr, carrier)

	assert.Equal(t, err.Error(), "malformed_token")
}

func TestTokenParse(t *testing.T) {
	tokenStr := Gen(Claims{"hi": "there"}, 0)
	token, _ := Parse(tokenStr)

	assert.Equal(t, "there", token.Claims.(jwt.MapClaims)["hi"].(string))

	tokenStr = "****"
	_, err := Parse(tokenStr)

	assert.Equal(t, "parse_token_failed", err.Error())
}

func TestTokenTest(t *testing.T) {
	tokenStr := Gen(Claims{"hi": "there"}, 0)
	err := Test(tokenStr)
	assert.NoError(t, err)

	tokenStr = "****"
	err = Test(tokenStr)
	assert.Equal(t, "parse_token_failed", err.Error())
}

func TestExpiredToken(t *testing.T) {
	tokenStr := Gen(Claims{"hi": "there"}, time.Duration(3))

	jwt.TimeFunc = func() time.Time {
		return time.Now().Local().Add(time.Hour * time.Duration(3))
	}

	err := Test(tokenStr)

	assert.NotEqual(t, err, nil)
	assert.Equal(t, "token_expired", err.Error())
}
