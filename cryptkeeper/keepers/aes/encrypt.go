package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

func encrypt(keyValue, unencrypted string) (Data, error) {
	result := Data{}
	key, _ := base64Decode([]byte(keyValue))
	plaintext := []byte(unencrypted)
	plaintext = pKCS5Padding(plaintext, aes.BlockSize)

	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ivtext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ivtext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return result, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(plaintext, plaintext)

	hash := make([]byte, hex.EncodedLen(len(plaintext)))
	hex.Encode(hash, plaintext)
	result.Hash = string(hash)

	salt := make([]byte, hex.EncodedLen(len(iv)))
	hex.Encode(salt, iv)
	result.Salt = string(salt)

	return result, nil
}

func base64Decode(message []byte) (b []byte, err error) {
	var l int
	b = make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	l, err = base64.StdEncoding.Decode(b, message)
	if err != nil {
		return
	}
	return b[:l], nil
}

func pKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}
