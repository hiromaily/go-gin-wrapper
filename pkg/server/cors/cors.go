package cors

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	str "github.com/hiromaily/go-gin-wrapper/pkg/strings"
)

const (
	// HeaderOrigin is for allowed specific domains
	HeaderOrigin = "Access-Control-Allow-Origin"
	// HeaderHeaders is for allowed headers from browser
	HeaderHeaders = "Access-Control-Allow-Headers"
	// HeaderMethods is for allowed method from browser
	HeaderMethods = "Access-Control-Allow-Methods"
	// HeaderCredentials is whether credentials are used of not
	HeaderCredentials = "Access-Control-Allow-Credentials"
)

// CORSer interface
type CORSer interface {
	ValidateHeader(ctx *gin.Context) error
	SetResponseHeader(ctx *gin.Context)
}

type cors struct {
	logger   *zap.Logger
	corsConf *config.CORS
}

// NewCORS returns CORSer interface
func NewCORS(
	logger *zap.Logger,
	corsConf *config.CORS,
) CORSer {
	return &cors{
		logger:   logger,
		corsConf: corsConf,
	}
}

// CheckHeader is for CORS before handling request
//  check preflight of XMLHttpRequest Level2 XMLHttpRequest
func (c *cors) ValidateHeader(ctx *gin.Context) error {
	c.logger.Info("cors CheckHeader")

	// TODO: should it be checked when `Options` method??
	// 1.check origin
	origin := ctx.Request.Header.Get("Origin")
	c.logger.Debug("CheckHeader", zap.String("origin header", origin))
	// Origin header: http://127.0.0.1:8000
	if str.SearchIndex(origin, c.corsConf.Origins) == -1 {
		return errors.New("origin header is invalid")
	}

	// 2.check header
	// added header intentionally set on Access-Control-Request-Headers
	// this value may change into lower case when extracting
	header := ctx.Request.Header.Get("Access-Control-Request-Headers")
	c.logger.Debug("CheckHeader", zap.String("Access-Control-Request-Headers header", header))
	// Access-Control-Request-Headers header: x-custom-header-cors, x-custom-header-gin

	// TODO: is it OK to check
	for _, h := range strings.Split(header, ",") {
		if str.SearchIndexLower(strings.TrimSpace(h), c.corsConf.Headers) == -1 {
			return errors.New("Access-Control-Request-Headers header is invalid")
		}
	}

	// 3.check method
	// this value is suppoused request to send after option method request
	method := ctx.Request.Header.Get("Access-Control-Request-Method")
	c.logger.Debug("CheckHeader", zap.String("Access-Control-Request-Method header", method))
	// Access-Control-Request-Method header: GET
	if str.SearchIndex(method, c.corsConf.Methods) == -1 {
		return errors.New("Access-Control-Request-Method header is invalid")
	}
	// if strings.Contains(s, "authorization") == true || strings.Contains(s, "Authorization") == true {

	return nil
}

// SetResponseHeader sets CORS header
// which has same type to gin.HandlerFunc `type HandlerFunc func(*Context)`
func (c *cors) SetResponseHeader(ctx *gin.Context) {
	c.logger.Info("cors SetHeader")
	if !c.corsConf.Enabled || ctx.Request.Method != "GET" {
		return
	}

	// Access-Control-Allow-Origin
	// allow from remote addr
	c.logger.Debug("", zap.String("c.Request.RemoteAddr", ctx.Request.RemoteAddr))
	// ctx.Writer.Header().Set(HeaderOrigin, ctx.Request.RemoteAddr)
	ctx.Writer.Header().Set(HeaderOrigin, strings.Join(c.corsConf.Origins, ", "))

	// Access-Control-Allow-Headers
	// ctx.Writer.Header().Set(HeaderHeaders, "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	ctx.Writer.Header().Set(HeaderHeaders, strings.Join(c.corsConf.Headers, ", "))

	// Access-Control-Allow-Methods
	// ctx.Writer.Header().Set(HeaderMethods, "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Writer.Header().Set(HeaderMethods, strings.Join(c.corsConf.Methods, ", "))

	// Access-Control-Allow-Credentials
	// ctx.Writer.Header().Set(HeaderCredentials, "false")
	ctx.Writer.Header().Set(HeaderCredentials, fmt.Sprintf("%t", c.corsConf.Credentials))
}
