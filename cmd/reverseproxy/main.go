package main

import (
	"flag"
	"log"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
)

var tomlPath = flag.String("f", "", "toml file path")

func init() {}

func main() {
	flag.Parse()

	// config
	conf, err := config.New(*tomlPath, false)
	if err != nil {
		panic(err)
	}

	// server
	regi := NewRegistry(conf)
	server := regi.NewProxyServer()
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
