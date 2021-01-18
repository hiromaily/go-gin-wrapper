package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hiromaily/go-gin-wrapper/pkg/encryption"
)

var (
	isEncoded = flag.Bool("encode", false, "encode target")
	isDecoded = flag.Bool("decode", false, "encode target")
)

var usage = `Usage: %s [options...]
Options:
  -encode
  -decode
e.g.:
  encryption -encode xxxxxxxx
    or
  encryption -decode xxxxxxxx
`

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
	}
	flag.Parse()

	if len(os.Args) != 3 {
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

	if *isEncoded {
		// encode
		log.Print(crypt.EncryptBase64(target))
		return
	}
	if *isDecoded {
		// decode
		decrypted, err := crypt.DecryptBase64(target)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(decrypted)
		return
	}
	flag.Usage()
}
