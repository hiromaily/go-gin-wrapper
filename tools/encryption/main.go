package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-gin-wrapper/pkg/encryption"
)

var mode = flag.String("m", "e", "e:encode, d:decode")

var usage = `Usage: %s [options...]
Options:
  -m  e:encode, d:decode.
e.g.:
  encryption -m e xxxxxxxx
    or
  encryption -m d xxxxxxxx
`

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
	}
	flag.Parse()

	if len(os.Args) != 4 {
		flag.Usage()
		os.Exit(1)
		return
	}
}

func setup() {
	key := os.Getenv("ENC_KEY")
	iv := os.Getenv("ENC_IV")

	if key == "" || iv == "" {
		log.Fatal(errors.New("set Environment Valuable: ENC_KEY, ENC_IV"))
		os.Exit(1)
	}

	encryption.NewCrypt(key, iv)
}

func main() {
	setup()

	crypt := encryption.GetCrypt()
	targetStr := os.Args[3]
	fmt.Printf("target string is %s\n", targetStr)

	switch *mode {
	case "e":
		// encode
		log.Print(crypt.EncryptBase64(targetStr))
	case "d":
		// decode
		str, err := crypt.DecryptBase64(targetStr)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(str)
	default:
		log.Fatal(errors.New("arguments is wrong"))
	}
}
