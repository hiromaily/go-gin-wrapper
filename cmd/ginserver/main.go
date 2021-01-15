package main

import (
	"flag"
	"log"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/encryption"
	"github.com/hiromaily/go-gin-wrapper/pkg/signal"
)

var (
	tomlPath        = flag.String("f", "", "toml file path")
	portNum         = flag.Int("p", 0, "port of server")
	isEncryptedConf = flag.Bool("crypto", false, "if true, config file is handled as encrypted value")
)

func init() {}

// Creates a gin router with default middleware:
// logger and recovery (crash-free) middleware
func main() {
	flag.Parse()

	// encryption
	if *isEncryptedConf {
		_, err := encryption.NewCryptWithEnv()
		if err != nil {
			panic(err)
		}
	}

	// config
	conf, err := config.NewInstance(*tomlPath, *isEncryptedConf)
	if err != nil {
		panic(err)
	}
	// FIXME: there are a lot of places singleton is used
	// config.New(*tomlPath, true)

	// debug mode
	if conf.Environment == "local" {
		// signal
		go signal.StartSignal()
	}

	isTestMode := false
	regi := NewRegistry(conf, isTestMode)
	server := regi.NewServer(*portNum)
	if _, err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
