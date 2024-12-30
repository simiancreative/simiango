package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

// Decrypt decrypts cipher text string into plain text string
func decrypt(keyValue, hash, salt string) (string, error) {
	txtBytes, _ := hex.DecodeString(hash)
	ivBytes, _ := hex.DecodeString(salt)
	key, _ := base64Decode([]byte(keyValue))

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(txtBytes)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(txtBytes, txtBytes)
	txtBytes = pKCS5UnPadding(txtBytes)

	return string(txtBytes), nil
}

func pKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
