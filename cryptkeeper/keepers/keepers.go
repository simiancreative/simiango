package keepers

import (
	"errors"
	"io"
)

type Configurator interface {
	Setup(input ...any) Config
}

type Encrypter interface {
	Decrypt(config Config, content io.Reader) (io.Reader, error)
	Encrypt(config Config, content io.Reader) (io.Reader, error)
}

type Hasher interface {
	Hash(config Config, password string) (string, error)
	Verify(config Config, hashedPassword, password string) bool
}

type Config map[string]string

func (c Config) Set(name, value string) Config {
	c[name] = value
	return c
}

func (c Config) Get(name string) string {
	val, ok := c[name]
	if !ok {
		return ""
	}

	return val
}

type Keeper struct {
	config       Config
	configurator Configurator
	encrypter    Encrypter
	hasher       Hasher
}

func (k *Keeper) SetConfigurator(c Configurator) *Keeper {
	k.configurator = c

	return k
}

func (k *Keeper) SetEncrypter(e Encrypter) *Keeper {
	k.encrypter = e

	return k
}

func (k *Keeper) SetHasher(h Hasher) *Keeper {
	k.hasher = h

	return k
}

func (k *Keeper) Setup(input ...any) *Keeper {
	if k.configurator == nil {
		return k
	}

	k.config = k.configurator.Setup(input...)

	return k
}

func (k *Keeper) Encrypt(content io.Reader) (io.Reader, error) {
	if k.encrypter == nil {
		return nil, errors.New("encrypter not implemented")
	}

	return k.encrypter.Encrypt(k.config, content)
}

func (k *Keeper) Decrypt(content io.Reader) (io.Reader, error) {
	if k.encrypter == nil {
		return nil, errors.New("encrypter not implemented")
	}

	return k.encrypter.Decrypt(k.config, content)
}

func (k *Keeper) Hash(password string) (string, error) {
	if k.hasher == nil {
		return "", errors.New("hasher not implemented")
	}

	return k.hasher.Hash(k.config, password)
}

func (k *Keeper) Verify(hashedPassword, password string) bool {
	if k.hasher == nil {
		return false
	}

	return k.hasher.Verify(k.config, hashedPassword, password)
}
