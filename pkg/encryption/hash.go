package encryption

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

//-----------------------------------------------------------------------------
// HASH
//-----------------------------------------------------------------------------

// GetMD5 is to hash by MD5
func GetMD5(baseString string) string {
	// md5
	h := md5.New()
	io.WriteString(h, baseString)
	ret := fmt.Sprintf("%x", h.Sum(nil))
	return ret
}

// GetSHA1 is to hash by SHA1
func GetSHA1(baseString string) string {
	// sha1
	h := sha1.New()
	io.WriteString(h, baseString)
	ret := fmt.Sprintf("%x", h.Sum(nil))
	return ret
}

// GetSHA256 is to hash by SHA256
func GetSHA256(baseString string) string {
	// sha256
	h := sha256.New()
	io.WriteString(h, baseString)
	ret := fmt.Sprintf("%x", h.Sum(nil))
	return ret
}

// GetMD5Plus is to hash by MD5 Plus salt
// 1.ユーザが入力したパスワードに対してMD5で一度暗号化
// 2.得られたMD5の値の前後に管理者自身だけが知っているランダムな文字列を追加
// 3.再度MD5で暗号化
func GetMD5Plus(baseString string, strPlus string) string {
	h := md5.New()
	io.WriteString(h, baseString)
	pwmd5 := fmt.Sprintf("%x", h.Sum(nil))

	salt1 := "@#$%"
	salt2 := "^&*()"

	// salt1 + username + salt2+MD5を連結。
	io.WriteString(h, salt1)
	io.WriteString(h, salt2)
	if strPlus != "" {
		io.WriteString(h, strPlus)
	}
	io.WriteString(h, pwmd5)

	ret := fmt.Sprintf("%x", h.Sum(nil))
	return ret
}

// GetScrypt to hash by Scrypt
func GetScrypt(baseString string) string {
	salt := "@#$%7G8r"
	// func Key(password, salt []byte, N, r, p, keyLen int) ([]byte, error) {
	dk, _ := scrypt.Key([]byte(baseString), []byte(salt), 16384, 8, 1, 32)

	// In order to read, it should be encoded by base64
	result := base64.StdEncoding.EncodeToString(dk)

	return result
}
