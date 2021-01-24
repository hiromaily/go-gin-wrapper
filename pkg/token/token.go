package token

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
)

// Generator interface
type Generator interface {
	Generate() string
}

type generator struct {
	salt string
}

// NewGenerator returns Generator
func NewGenerator(salt string) Generator {
	return &generator{
		salt: salt,
	}
}

func (g *generator) Generate() string {
	md5Hash := md5.New()

	io.WriteString(md5Hash, strconv.FormatInt(time.Now().UnixNano(), 10))
	io.WriteString(md5Hash, g.salt)
	return fmt.Sprintf("%x", md5Hash.Sum(nil))
}
