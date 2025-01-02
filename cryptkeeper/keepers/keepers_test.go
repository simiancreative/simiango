package keepers_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/simiancreative/simiango/cryptkeeper/keepers"
	"github.com/stretchr/testify/assert"
)

type mockConfigurator struct{}

func (m *mockConfigurator) Setup(input ...any) keepers.Config {
	return keepers.Config{"key": "value"}
}

type mockEncrypter struct{}

func (m *mockEncrypter) Encrypt(config keepers.Config, content io.Reader) (io.Reader, error) {
	return bytes.NewBufferString("encrypted"), nil
}

func (m *mockEncrypter) Decrypt(config keepers.Config, content io.Reader) (io.Reader, error) {
	return bytes.NewBufferString("decrypted"), nil
}

type mockHasher struct{}

func (m *mockHasher) Hash(config keepers.Config, password string) (string, error) {
	return "hashedPassword", nil
}

func (m *mockHasher) Verify(config keepers.Config, hashedPassword, password string) bool {
	return hashedPassword == "hashedPassword" && password == "password"
}

func TestKeeper_SetConfigurator(t *testing.T) {
	k := &keepers.Keeper{}
	k.SetConfigurator(&mockConfigurator{})
	assert.NotNil(t, k)
	assert.NotNil(t, k.SetConfigurator(&mockConfigurator{}))
}

func TestKeeper_SetEncrypter(t *testing.T) {
	k := &keepers.Keeper{}
	k.SetEncrypter(&mockEncrypter{})
	assert.NotNil(t, k)
	assert.NotNil(t, k.SetEncrypter(&mockEncrypter{}))
}

func TestKeeper_SetHasher(t *testing.T) {
	k := &keepers.Keeper{}
	k.SetHasher(&mockHasher{})
	assert.NotNil(t, k)
	assert.NotNil(t, k.SetHasher(&mockHasher{}))
}

func TestKeeper_Encrypt(t *testing.T) {
	k := &keepers.Keeper{}
	k.SetEncrypter(&mockEncrypter{})
	result, err := k.Encrypt(bytes.NewBufferString("content"))
	assert.NoError(t, err)
	buf := new(bytes.Buffer)
	buf.ReadFrom(result)
	assert.Equal(t, "encrypted", buf.String())
}

func TestKeeper_Decrypt(t *testing.T) {
	k := &keepers.Keeper{}
	k.SetEncrypter(&mockEncrypter{})
	result, err := k.Decrypt(bytes.NewBufferString("content"))
	assert.NoError(t, err)
	buf := new(bytes.Buffer)
	buf.ReadFrom(result)
	assert.Equal(t, "decrypted", buf.String())
}

func TestKeeper_Hash(t *testing.T) {
	k := &keepers.Keeper{}
	k.SetHasher(&mockHasher{})
	result, err := k.Hash("password")
	assert.NoError(t, err)
	assert.Equal(t, "hashedPassword", result)
}

func TestKeeper_Verify(t *testing.T) {
	k := &keepers.Keeper{}
	k.SetHasher(&mockHasher{})
	result := k.Verify("hashedPassword", "password")
	assert.True(t, result)
}
