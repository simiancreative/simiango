package token

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/simiancreative/simiango/context"
)

func Test(token string) error {
	t, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(getKey()), nil
		},
	)

	if err != nil {
		return errors.New("token_invalid")
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

	claims["jti"] = context.Id()
	claims["iat"] = time.Now().Unix()

	for k, v := range params {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(getKey())

	return tokenString
}

func getKey() []byte {
	key := os.Getenv("TOKEN_SECRET")
	return []byte(key)
}
