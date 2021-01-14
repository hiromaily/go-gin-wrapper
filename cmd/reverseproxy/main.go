package main

import (
	"flag"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	lg "github.com/hiromaily/golibs/log"
)

var tomlPath = flag.String("f", "", "Toml file path")

func init() {}

func main() {
	flag.Parse()

	conf, err := config.NewInstance(*tomlPath, false)
	if err != nil {
		panic(err)
	}

	// log
	logLevel := lg.LogStatus(conf.Proxy.Server.Log.Level)
	lg.InitializeLog(logLevel, lg.TimeShortFile,
		"[REVERSE_PROXY]", conf.Proxy.Server.Log.Path, "hiromaily")

	regi := NewRegistry(conf)
	server := regi.NewProxyServerer()
	if err := server.Start(); err != nil {
		lg.Error(err)
	}
}
