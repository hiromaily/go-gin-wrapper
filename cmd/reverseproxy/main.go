package main

import (
	"flag"
	"log"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
)

var tomlPath = flag.String("f", "", "Toml file path")

func init() {}

func main() {
	flag.Parse()

	conf, err := config.NewInstance(*tomlPath, false)
	if err != nil {
		panic(err)
	}

	regi := NewRegistry(conf)
	server := regi.NewProxyServerer()
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
