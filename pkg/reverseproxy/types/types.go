package types

// ProxyMode is proxy type
type ProxyMode string

const (
	// NoProxy doesn't use prooxy
	NoProxy ProxyMode = "no"
	// GoGinProxy uses ./cmd/reverseproxy as proxy server
	GoGinProxy ProxyMode = "go-gin-proxy"
	// NginxProxy uses Nginx as proxy server
	NginxProxy ProxyMode = "nginx"
)
