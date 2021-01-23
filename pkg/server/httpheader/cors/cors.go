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

// ValidateHeader validates CORS before handling request
//  check preflight (OPTIONS) of XMLHttpRequest Level2 XMLHttpRequest
func (c *cors) ValidateHeader(ctx *gin.Context) error {
	origin := ctx.Request.Header.Get("Origin")
	aclHeader := ctx.Request.Header.Get("Access-Control-Request-Headers")
	method := ctx.Request.Header.Get("Access-Control-Request-Method")

	c.logger.Info("cors CheckHeader",
		zap.String("origin header", origin),
		zap.String("Access-Control-Request-Headers header", aclHeader),
		zap.String("Access-Control-Request-Method header", method),
	)

	// Origin header would be `http://127.0.0.1:8000` on local
	if str.SearchIndex(origin, c.corsConf.Origins) == -1 {
		return errors.New("origin header is invalid")
	}

	// Access-Control-Request-Headers
	for _, header := range strings.Split(aclHeader, ",") {
		if str.SearchIndexLower(strings.TrimSpace(header), c.corsConf.Headers) == -1 {
			return errors.Errorf("Access-Control-Request-Headers: %s is invalid", header)
		}
	}

	// Access-Control-Request-Method header: GET
	if str.SearchIndex(method, c.corsConf.Methods) == -1 {
		return errors.New("Access-Control-Request-Method header is invalid")
	}
	return nil
}

// SetResponseHeader sets CORS header
// which has same type to gin.HandlerFunc `type HandlerFunc func(*Context)`
func (c *cors) SetResponseHeader(ctx *gin.Context) {
	c.logger.Info("cors SetHeader",
		zap.String("c.Request.RemoteAddr", ctx.Request.RemoteAddr),
	)
	if !c.corsConf.Enabled || ctx.Request.Method != "GET" {
		return
	}

	// Access-Control-Allow-Origin: allow from remote addr
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
