package reverseproxy

import (
	"container/ring"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	lg "github.com/hiromaily/golibs/log"
)

// Serverer is Serverer interface
type Serverer interface {
	Start() error
}

// NewServerer is to return Serverer interface
func NewServerer(conf *config.Config) Serverer {
	return NewServer(conf)
}

// ----------------------------------------------------------------------------
// Server
// ----------------------------------------------------------------------------

// Server is Server object
type Server struct {
	conf *config.Config
}

// NewServer is to return server object
func NewServer(
	conf *config.Config) *Server {
	srv := Server{
		conf: conf,
	}
	return &srv
}

// Start is to start server execution
func (s *Server) Start() error {
	ports := s.conf.Proxy.Server.WebPort

	if len(ports) == 1 {
		s.singleReverseProxy()
	} else if len(ports) > 1 {
		s.multipleReverseProxy()
	}

	return nil
}

// Single Reverse Proxy
func (s *Server) singleReverseProxy() {
	lg.Info("singleReverseProxy()")
	// Web Server
	// webserverURL := "http://127.0.0.1:9990"
	srv := s.conf.Server
	tmp := getURL(srv.Scheme, srv.Host, srv.Port)
	webserverURL, _ := url.Parse(tmp)

	fmt.Printf("proxy is runnig ... using Port: %d\n", s.conf.Proxy.Server.Port)

	// Proxy Server
	proxyAddress := fmt.Sprintf(":%d", s.conf.Proxy.Server.Port)
	proxyHandler := httputil.NewSingleHostReverseProxy(webserverURL)
	server := http.Server{
		Addr:    proxyAddress,
		Handler: proxyHandler,
	}
	server.ListenAndServe()
}

// Multiple Reverse Proxy
func (s *Server) multipleReverseProxy() {
	ports := s.conf.Proxy.Server.WebPort
	lg.Infof("multipleReverseProxy(): number of servers is %d", len(ports))
	// As precondition, increment port number by one.

	// web servers
	srv := s.conf.Server
	hostRing := ring.New(len(ports))
	for _, port := range ports {
		// url, _ := url.Parse(getURL(srv.Scheme, srv.Host, srv.Port+i))
		url, _ := url.Parse(getURL(srv.Scheme, srv.Host, port))
		hostRing.Value = url
		hostRing = hostRing.Next()
	}

	mutex := sync.Mutex{}
	// access server alternately
	director := func(request *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()
		request.URL.Scheme = srv.Scheme
		request.URL.Host = hostRing.Value.(*url.URL).Host
		hostRing = hostRing.Next()
	}

	// proxy
	proxy := &httputil.ReverseProxy{Director: director}
	proxyAddress := fmt.Sprintf(":%d", s.conf.Proxy.Server.Port)

	server := http.Server{
		Addr:    proxyAddress,
		Handler: proxy,
	}
	server.ListenAndServe()
}

func getURL(scheme, host string, port int) string {
	return fmt.Sprintf("%s://%s:%d", scheme, host, port)
}
