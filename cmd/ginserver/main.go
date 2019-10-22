package main

import (
	"flag"

	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	lg "github.com/hiromaily/golibs/log"
	"github.com/hiromaily/golibs/signal"
	"github.com/hiromaily/golibs/cipher/encryption"
)

var (
	tomlPath = flag.String("f", "", "toml file path")
	portNum  = flag.Int("p", 0, "port of server")
	isEncryptedConf = flag.Bool("crypto", false, "if true, config file is handled as encrypted value")
)

func init() {}

func parseFlag() {
	//command-line
	flag.Parse()
}

// Creates a gin router with default middleware:
// logger and recovery (crash-free) middleware
func main() {
	parseFlag()

	//cipher
	if *isEncryptedConf {
		_, err := encryption.NewCryptWithEnv()
		if err != nil {
			panic(err)
		}
	}

	//config
	conf, err := configs.NewInstance(*tomlPath, *isEncryptedConf)
	if err != nil {
		panic(err)
	}
	//FIXME: there are a lot of places singleton is used
	configs.New(*tomlPath, true)

	//log
	lg.InitializeLog(lg.LogStatus(conf.Server.Log.Level), lg.TimeShortFile,
		"[GOWEB]", conf.Server.Log.Path, "hiromaily")

	// debug mode
	if conf.Environment == "local" {
		//signal
		go signal.StartSignal()
	}

	regi := NewRegistry(conf)
	server := regi.NewServerer(*portNum)
	if err := server.Start(); err != nil {
		lg.Error(err)
	}
}
