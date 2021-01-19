package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"

	"github.com/pkg/errors"
)

// https://github.com/tadzik/simpleaes/blob/master/simpleaes.go

// Crypt interface
type Crypt interface {
	Encrypt(src []byte) []byte
	Decrypt(src []byte) []byte
	EncryptBase64(plainText string) string
	DecryptBase64(base64String string) (string, error)
}

type crypt struct {
	cipher cipher.Block
	iv     []byte
}

// Creates a new encryption/decryption object
// with a given key of a given size
// (16, 24 or 32 for AES-128, AES-192 and AES-256 respectively,
// as per http://golang.org/pkg/crypto/aes/#NewCipher)
//
// The key will be padded to the given size if needed.
// An IV is created as a series of NULL bytes of necessary length
// when there is no iv string passed as 3rd value to function.

// NewCrypt returns Crypt interface
// key size should be 16,24,32
// iv size should be 16
func NewCrypt(key, iv string) (Crypt, error) {
	if key == "" || iv == "" {
		return nil, errors.New("both of key and iv is required")
	}

	padded := make([]byte, len(key))
	copy(padded, []byte(key))

	bIv := []byte(iv)
	block, err := aes.NewCipher(padded)
	if err != nil {
		return nil, err
	}

	return &crypt{block, bIv}, nil
}

// NewCryptWithEnv returns Crypt interface created by environment variable.
func NewCryptWithEnv() (Crypt, error) {
	key := os.Getenv("ENC_KEY")
	iv := os.Getenv("ENC_IV")
	if key == "" || iv == "" {
		return nil, errors.Errorf("%s", "Environment Variable: `ENC_KEY`, `ENC_IV` is required")
	}

	return NewCrypt(key, iv)
}

func (c *crypt) padSlice(src []byte) []byte {
	// src must be a multiple of block size
	mult := (len(src) / aes.BlockSize) + 1
	leng := aes.BlockSize * mult

	srcPadded := make([]byte, leng)
	copy(srcPadded, src)
	return srcPadded
}

// Encrypt encrypts a slice of bytes, producing a new, freshly allocated slice
// Source will be padded with null bytes if necessary
func (c *crypt) Encrypt(src []byte) []byte {
	if len(src)%aes.BlockSize != 0 {
		src = c.padSlice(src)
	}
	dst := make([]byte, len(src))
	cipher.NewCBCEncrypter(c.cipher, c.iv).CryptBlocks(dst, src)
	return dst
}

// Decrypt decrypts a slice of bytes, producing a new, freshly allocated slice
// Source will be padded with null bytes if necessary
func (c *crypt) Decrypt(src []byte) []byte {
	if len(src)%aes.BlockSize != 0 {
		src = c.padSlice(src)
	}
	dst := make([]byte, len(src))
	cipher.NewCBCDecrypter(c.cipher, c.iv).CryptBlocks(dst, src)
	trimmed := bytes.Trim(dst, "\x00")
	return trimmed
}

// EncryptBase64 encrypts and encode by base64
func (c *crypt) EncryptBase64(plainText string) string {
	encryptedBytes := c.Encrypt([]byte(plainText))
	base64 := base64.StdEncoding.EncodeToString(encryptedBytes)
	return base64
}

// DecryptBase64 decrypts decoded Base64 string
func (c *crypt) DecryptBase64(base64String string) (string, error) {
	unbase64, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", err
	}
	decryptedBytes := c.Decrypt(unbase64)
	return string(decryptedBytes[:]), nil
}
