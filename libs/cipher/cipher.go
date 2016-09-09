package cipher

import (
	enc "github.com/hiromaily/golibs/cipher/encryption"
	"os"
)

// Setup is for setup
func Setup() {
	size := 16
	key := os.Getenv("ENC_KEY")
	iv := os.Getenv("ENC_IV")

	if key == "" || iv == "" {
		panic("set Environment Variable: ENC_KEY, ENC_IV")
	}

	enc.NewCrypt(size, key, iv)
}

//crypt := enc.GetCryptInstance()

//encode
//crypt.EncryptBase64(targetStr)

//decode
//crypt.DecryptBase64(targetStr)
