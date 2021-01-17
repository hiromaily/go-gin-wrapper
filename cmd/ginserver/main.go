package main

import (
	"flag"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/encryption"
)

var (
	tomlPath    = flag.String("f", "", "toml file path")
	portNum     = flag.Int("p", 0, "sever port")
	isEncrypted = flag.Bool("crypto", false, "if true, values in config file are encrypted")
)

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
	conf, err := config.NewConfig(*tomlPath, *isEncrypted)
	if err != nil {
		panic(err)
	}

	// overwrite config by args
	if *portNum != 0 {
		conf.Server.Port = *portNum
	}

	// server
	regi := NewRegistry(conf, false)
	server := regi.NewServer()
	if _, err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
