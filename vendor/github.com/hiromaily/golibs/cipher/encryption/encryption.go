package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	cph "crypto/cipher"
	"encoding/base64"
	"io"
)

//https://github.com/tadzik/simpleaes/blob/master/simpleaes.go

// Crypt is for cipher config data
type Crypt struct {
	//enc, dec cipher.BlockMode
	cipher cipher.Block
	iv     []byte
}

var (
	cryptInfo Crypt
	commonIV  = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
)

// Creates a new encryption/decryption object
// with a given key of a given size
// (16, 24 or 32 for AES-128, AES-192 and AES-256 respectively,
// as per http://golang.org/pkg/crypto/aes/#NewCipher)
//
// The key will be padded to the given size if needed.
// An IV is created as a series of NULL bytes of necessary length
// when there is no iv string passed as 3rd value to function.
//func NewCryptUtil(size int, key string, more ...string) (*CryptUtil, error) {

// NewCrypt is to create crypt instance
// e.g. size = 16, key = "8#75F%R+&a5ZvM_<", iv = "@~wp-7hPs<WEx@R4"
func NewCrypt(size int, key, iv string) (*Crypt, error) {

	//fmt.Printf("cryptConf %#v\n", cryptConf)

	padded := make([]byte, size)
	copy(padded, []byte(key))

	bIv := []byte(iv)
	aes, err := aes.NewCipher(padded)
	if err != nil {
		return nil, err
	}

	cryptInfo = Crypt{aes, bIv}

	return &cryptInfo, nil
}

// GetCrypt is to get crypt instance
func GetCrypt() *Crypt {
	return &cryptInfo
}

func (c *Crypt) padSlice(src []byte) []byte {
	// src must be a multiple of block size
	bs := 16
	mult := int((len(src) / bs) + 1)
	leng := bs * mult

	srcPadded := make([]byte, leng)
	copy(srcPadded, src)
	return srcPadded
}

// Encrypt is encrypt a slice of bytes, producing a new, freshly allocated slice
// Source will be padded with null bytes if necessary
func (c *Crypt) Encrypt(src []byte) []byte {
	if len(src)%16 != 0 {
		src = c.padSlice(src)
	}
	dst := make([]byte, len(src))
	cipher.NewCBCEncrypter(c.cipher, c.iv).CryptBlocks(dst, src)
	return dst
}

// EncryptBase64 is encrypt and encode by base64 string
func (c *Crypt) EncryptBase64(plainText string) string {
	encryptedBytes := c.Encrypt([]byte(plainText))
	base64 := base64.StdEncoding.EncodeToString(encryptedBytes)
	return base64
}

// EncryptStream is to encrypt blocks from reader, write results into writer
func (c *Crypt) EncryptStream(reader io.Reader, writer io.Writer) error {
	for {
		buf := make([]byte, 16)
		_, err := io.ReadFull(reader, buf)
		if err != nil {
			if err == io.EOF {
				break
			} else if err == io.ErrUnexpectedEOF {
				// nothing
			} else {
				return err
			}
		}
		cipher.NewCBCEncrypter(c.cipher, c.iv).CryptBlocks(buf, buf)
		if _, err = writer.Write(buf); err != nil {
			return err
		}
	}
	return nil
}

// Decrypt is to decrypt a slice of bytes, producing a new, freshly allocated slice
// Source will be padded with null bytes if necessary
func (c *Crypt) Decrypt(src []byte) []byte {
	if len(src)%16 != 0 {
		src = c.padSlice(src)
	}
	dst := make([]byte, len(src))
	cipher.NewCBCDecrypter(c.cipher, c.iv).CryptBlocks(dst, src)
	trimmed := bytes.Trim(dst, "\x00")
	return trimmed
}

// DecryptBase64 is to decrypt decoded Base64 string
func (c *Crypt) DecryptBase64(base64String string) (string, error) {
	unbase64, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", err
	}
	decryptedBytes := c.Decrypt(unbase64)
	return string(decryptedBytes[:]), nil
}

// DecryptStream is to decrypt blocks from reader, write results into writer
func (c *Crypt) DecryptStream(reader io.Reader, writer io.Writer) error {
	buf := make([]byte, 16)
	for {
		_, err := io.ReadFull(reader, buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		cipher.NewCBCDecrypter(c.cipher, c.iv).CryptBlocks(buf, buf)
		if _, err = writer.Write(buf); err != nil {
			return err
		}
	}
	return nil
}

//-----------------------------------------------------------------------------
// Cipher (TODO:It hasn't finished yet)
//-----------------------------------------------------------------------------

// GetAesEncrypt to cipher by Aes
func GetAesEncrypt(baseString string) (string, error) {
	//The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
	key := "djjf63Hdgd#:dj37"

	// create　aes cipher algorithm
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// 16 bytes for AES-128, 24 bytes for AES-192, 32 bytes for AES-256
	cipherText := []byte("akcgey87275r78jp")
	iv := cipherText[:aes.BlockSize] // const BlockSize = 16

	targetStr := []byte(baseString)

	// encrypt
	encrypter := cph.NewCFBEncrypter(block, iv)
	encrypted := make([]byte, len(targetStr))
	encrypter.XORKeyStream(encrypted, targetStr)

	// decrypt
	//TODO:通常このタイミングで元のサイズってわからないのでは？
	decrypter := cph.NewCFBDecrypter(block, iv)
	decrypted := make([]byte, len(targetStr))
	decrypter.XORKeyStream(decrypted, encrypted)

	return "", nil
}

//-----------------------------------------------------------------------------
// Base64
//-----------------------------------------------------------------------------

// GetBase64Encode is to encode by Base64
func GetBase64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

// GetBase64Decode is to decode by base64
func GetBase64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}
