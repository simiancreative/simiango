package cryptkeeper

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"os"
)

var keyName = "TOKEN_SECRET"

func SetKeyName(newName string) {
	keyName = newName
}

// Decrypt decrypts cipher text string into plain text string
func Decrypt(hash string, salt string) (string, error) {
	txtBytes, _ := hex.DecodeString(hash)
	ivBytes, _ := hex.DecodeString(salt)
	key, _ := Base64Decode([]byte(os.Getenv(keyName)))

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(txtBytes)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(txtBytes, txtBytes)
	txtBytes = PKCS5UnPadding(txtBytes)

	return string(txtBytes), nil
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
