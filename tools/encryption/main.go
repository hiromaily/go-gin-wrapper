package main

import (
	"flag"
	"fmt"
	"os"

	enc "github.com/hiromaily/golibs/cipher/encryption"
	lg "github.com/hiromaily/golibs/log"
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
	lg.InitializeLog(lg.DebugStatus, lg.TimeShortFile, "[GOTOOLS GoChipher]", "", "hiromaily")

	key := os.Getenv("ENC_KEY")
	iv := os.Getenv("ENC_IV")

	if key == "" || iv == "" {
		lg.Fatalf("%s", "set Environment Valuable: ENC_KEY, ENC_IV")
		os.Exit(1)
	}

	enc.NewCrypt(key, iv)
}

func main() {
	setup()

	crypt := enc.GetCrypt()
	targetStr := os.Args[3]
	fmt.Printf("target string is %s\n", targetStr)

	switch *mode {
	case "e":
		// encode
		lg.Info(crypt.EncryptBase64(targetStr))
	case "d":
		// decode
		str, err := crypt.DecryptBase64(targetStr)
		if err != nil {
			lg.Fatal(err)
		}
		lg.Info(str)
	default:
		lg.Fatalf("%s", "arguments is wrong")
	}
}
