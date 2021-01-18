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

func main() {
	crypt, err := encryption.NewCryptWithEnv()
	if err != nil {
		panic(err)
	}

	target := os.Args[3]
	log.Printf("target string is %s\n", target)

	switch *mode {
	case "e":
		// encode
		log.Print(crypt.EncryptBase64(target))
	case "d":
		// decode
		decrypted, err := crypt.DecryptBase64(target)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(decrypted)
	default:
		flag.Usage()
		log.Fatal(errors.New("arguments `-m` is wrong"))
	}
}
