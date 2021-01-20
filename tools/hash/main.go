package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hiromaily/go-gin-wrapper/pkg/encryption"
)

var (
	isMD5  = flag.Bool("md5", false, "md5 algo")
	salt1  = flag.String("salt1", "", "salt1")
	salt2  = flag.String("salt2", "", "salt2")
	target = flag.String("target", "", "target string")
)

var usage = `Usage: %s [options...]
Options:
  -md5
  -salt1
  -salt2
  -target
e.g.:
  hash -md5 -salt1 xxx -salt2 xxx -target xxx
`

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
	}
	flag.Parse()

	if len(os.Args) != 8 {
		flag.Usage()
		os.Exit(1)
		return
	}
}

func main() {
	if *isMD5 && *salt1 != "" && *salt2 != "" && *target != "" {
		log.Print("hash by md5")
		hashMD5 := encryption.NewMD5(*salt1, *salt2)
		log.Print(hashMD5.Hash(*target))
		return
	}
	flag.Usage()
}
