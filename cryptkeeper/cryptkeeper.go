package cryptkeeper

import (
	"github.com/simiancreative/simiango/cryptkeeper/keepers"

	"github.com/simiancreative/simiango/cryptkeeper/keepers/aes"
	"github.com/simiancreative/simiango/cryptkeeper/keepers/bcrypt"
	"github.com/simiancreative/simiango/cryptkeeper/keepers/pgp"
)

const (
	AES = iota
	PGP
	BCRYPT
)

var names = map[int]string{
	AES:    "AES",
	PGP:    "PGP",
	BCRYPT: "BCRYPT",
}

func init() {
	register(AES, aes.New)
	register(PGP, pgp.New)
	register(BCRYPT, bcrypt.New)
}

var registered = map[int]func() *keepers.Keeper{}

func register(n int, k func() *keepers.Keeper) {
	registered[n] = k
}

func New(kind int) (*keepers.Keeper, error) {
	keeper, ok := registered[kind]
	if !ok {
		return nil, NotFoundError(kind)
	}

	return keeper(), nil
}
