package main

import (
	"flag"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/encryption"
	"github.com/hiromaily/go-gin-wrapper/pkg/signal"
)

var (
	tomlPath    = flag.String("f", "", "toml file path")
	portNum     = flag.Int("p", 0, "port of server")
	isEncrypted = flag.Bool("crypto", false, "if true, values in config file are encrypted")
)

func init() {}

func main() {
	flag.Parse()

	// encryption
	if *isEncrypted {
		_, err := encryption.NewCryptWithEnv()
		if err != nil {
			panic(err)
		}
	}

	// config
	conf, err := config.New(*tomlPath, *isEncrypted)
	if err != nil {
		panic(err)
	}

	// overwrite config by args
	if *portNum != 0 {
		conf.Server.Port = *portNum
	}

	// accept signal
	if conf.IsSignal {
		go signal.StartSignal()
	}

	regi := NewRegistry(conf, false)
	server := regi.NewServer()
	if _, err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
