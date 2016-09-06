package main

import (
	"container/ring"
	"flag"
	"fmt"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var (
	tomlPath = flag.String("f", "", "Toml file path")
)

func init() {
	//command-line
	flag.Parse()

	//conf
	initConf()
}

func initConf() {
	//config
	if *tomlPath != "" {
		conf.SetTomlPath(*tomlPath)
	}
	conf.New("")

	//log
	lg.InitializeLog(conf.GetConf().Proxy.Server.Log.Level, lg.LOG_OFF_COUNT, 0,
		"[REVERSE_PROXY]", conf.GetConf().Proxy.Server.Log.Path)
}

func getURL(scheme, host string, port int) string {
	return fmt.Sprintf("%s://%s:%d", scheme, host, port)
}

func singleReverseProxy() {
	lg.Info("singleReverseProxy()")
	//Web Server
	//webserverURL := "http://127.0.0.1:9990"
	srv := conf.GetConf().Server
	webserverURL := getURL(srv.Scheme, srv.Host, srv.Port)
	webserverUrl, _ := url.Parse(webserverURL)

	//Proxy Server
	proxyAddress := fmt.Sprintf(":%d", conf.GetConf().Proxy.Server.Port)
	proxyHandler := httputil.NewSingleHostReverseProxy(webserverUrl)
	server := http.Server{
		Addr:    proxyAddress,
		Handler: proxyHandler,
	}
	server.ListenAndServe()
}

func multipleReverseProxy(num int) {
	lg.Info("multipleReverseProxy() number of servers is %d", num)
	//As precondition, increment port number by one.

	//web servers
	srv := conf.GetConf().Server
	hostRing := ring.New(num)
	for i := 0; i < num; i++ {
		url, _ := url.Parse(getURL(srv.Scheme, srv.Host, srv.Port+i))
		hostRing.Value = url
		hostRing = hostRing.Next()
	}

	mutex := sync.Mutex{}
	director := func(request *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()
		request.URL.Scheme = srv.Scheme
		request.URL.Host = hostRing.Value.(*url.URL).Host
		hostRing = hostRing.Next()
	}

	//proxy
	proxy := &httputil.ReverseProxy{Director: director}
	proxyAddress := fmt.Sprintf(":%d", conf.GetConf().Proxy.Server.Port)

	server := http.Server{
		Addr:    proxyAddress,
		Handler: proxy,
	}
	server.ListenAndServe()
}

func main() {
	//TODO:define on outer file.
	serverNum := 1

	if serverNum == 1 {
		singleReverseProxy()
	} else if serverNum > 1 {
		multipleReverseProxy(serverNum)
	}
}
