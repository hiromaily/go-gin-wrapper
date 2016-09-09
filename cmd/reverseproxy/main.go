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
	ports    []int
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
		conf.SetTOMLPath(*tomlPath)
	}
	conf.New("")

	//log
	lg.InitializeLog(conf.GetConf().Proxy.Server.Log.Level, lg.LOG_OFF_COUNT, 0,
		"[REVERSE_PROXY]", conf.GetConf().Proxy.Server.Log.Path)
}

func getURL(scheme, host string, port int) string {
	return fmt.Sprintf("%s://%s:%d", scheme, host, port)
}

//Single Reverse Proxy
func singleReverseProxy() {
	lg.Info("singleReverseProxy()")
	//Web Server
	//webserverURL := "http://127.0.0.1:9990"
	srv := conf.GetConf().Server
	tmp := getURL(srv.Scheme, srv.Host, srv.Port)
	webserverURL, _ := url.Parse(tmp)

	//Proxy Server
	proxyAddress := fmt.Sprintf(":%d", conf.GetConf().Proxy.Server.Port)
	proxyHandler := httputil.NewSingleHostReverseProxy(webserverURL)
	server := http.Server{
		Addr:    proxyAddress,
		Handler: proxyHandler,
	}
	server.ListenAndServe()
}

// Multiple Reverse Proxy
func multipleReverseProxy() {
	lg.Infof("multipleReverseProxy(): number of servers is %d", len(ports))
	//As precondition, increment port number by one.

	//web servers
	srv := conf.GetConf().Server
	hostRing := ring.New(len(ports))
	for _, port := range ports {
		//url, _ := url.Parse(getURL(srv.Scheme, srv.Host, srv.Port+i))
		url, _ := url.Parse(getURL(srv.Scheme, srv.Host, port))
		hostRing.Value = url
		hostRing = hostRing.Next()
	}

	mutex := sync.Mutex{}
	//access server alternately
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
	ports = conf.GetConf().Proxy.Server.WebPort

	fmt.Printf("proxy is runnig ... using Port: %d\n", conf.GetConf().Proxy.Server.Port)

	if len(ports) == 1 {
		singleReverseProxy()
	} else if len(ports) > 1 {
		multipleReverseProxy()
	}

}
