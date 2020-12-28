package token

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/simiancreative/simiango/meta"
)

// Decode accepts a token string and unmarshals the header and payload
// segments into the given interface{}
func Decode(token string, v interface{}) error {
	parts := strings.Split(token, ".")

	decoded, err := jwt.DecodeSegment(parts[0])
	if err != nil {
		return err
	}
	json.Unmarshal(decoded, v)

	decoded, err = jwt.DecodeSegment(parts[1])
	if err != nil {
		return err
	}
	json.Unmarshal(decoded, v)

	return nil
}

// ParseWithSecret attempts to parse a token string given a custom siginin key.
// Then returns a Token if string/key is valid
func ParseWithSecret(token string, secret []byte) (*jwt.Token, error) {
	t, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(getKey()), nil
		},
	)

	if err != nil {
		return nil, errors.New("parse_token_failed")
	}

	return t, nil
}

// Parse does what ParseWithSecret does but gets the signing key from the env
// var TOKEN_SECRET
func Parse(token string) (*jwt.Token, error) {
	return ParseWithSecret(token, getKey())
}

// TestWithSecret verifies a token and returns an error if the token is invalid
func TestWithSecret(token string, secret []byte) error {
	t, err := Parse(token)

	if err != nil {
		return err
	}

	if !t.Valid {
		return errors.New("token_invalid")
	}

	return nil
}

// Test does what TestWithSecret does but gets the signing key from the env
// var TOKEN_SECRET
func Test(token string) error {
	return TestWithSecret(token, getKey())
}

// Claims defines the structure for a token payload
type Claims map[string]interface{}

// GenWithSecret generates a jwt token string from a payload, secret and exp
// duration
func GenWithSecret(params Claims, secret []byte, expMinutes time.Duration) string {
	claims := jwt.MapClaims{}

	if expMinutes > 0 {
		claims["exp"] = time.Now().Add(time.Minute * expMinutes).Unix()
	}

	claims["jti"] = meta.Id()
	claims["iat"] = time.Now().Unix()

	for k, v := range params {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(secret)

	return tokenString
}

// Gen does what GenWithSecret does but gets the signing key from the env
// var TOKEN_SECRET
func Gen(params Claims, expMinutes time.Duration) string {
	return GenWithSecret(params, getKey(), expMinutes)
}

func getKey() []byte {
	key, ok := os.LookupEnv("TOKEN_SECRET")
	if !ok {
		panic(errors.New("token_signing_key_not_configured"))
	}

	return []byte(key)
}
