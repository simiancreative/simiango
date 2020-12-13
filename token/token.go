package token

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/simiancreative/simiango/meta"
)

func Parse(token string) (*jwt.Token, error) {
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

func Test(token string) error {
	t, err := Parse(token)

	if err != nil {
		return err
	}

	if !t.Valid {
		return errors.New("token_invalid")
	}

	return nil
}

type Claims map[string]interface{}

func Gen(params Claims, expMinutes time.Duration) string {
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
	tokenString, _ := token.SignedString(getKey())

	return tokenString
}

func getKey() []byte {
	key, ok := os.LookupEnv("TOKEN_SECRET")
	if !ok {
		panic(errors.New("token_signing_key_not_configured"))
	}

	return []byte(key)
}
