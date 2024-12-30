package bcrypt

import (
	"github.com/simiancreative/simiango/cryptkeeper/keepers"
	"golang.org/x/crypto/bcrypt"
)

type Keeper struct{}

func (p Keeper) Verify(_ keepers.Config, hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (p Keeper) Hash(_ keepers.Config, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func New() *keepers.Keeper {
	p := Keeper{}
	k := &keepers.Keeper{}

	return k.
		SetHasher(p)
}
