package ginctx

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/server/ginurl"
)

// DebugContext debugs context
func DebugContext(ctx *gin.Context, logger *zap.Logger) {
	logger.Debug("request",
		zap.Any("gin_ctx", ctx),
		zap.Any("gin_ctx_keys", ctx.Keys),
		zap.String("request_method", ctx.Request.Method),
		zap.Any("request_header", ctx.Request.Header),
		zap.Any("request_body", ctx.Request.Body),
		zap.Any("request_url", ctx.Request.URL),
		zap.String("request_url_string", ginurl.GetURLString(ctx)),
		zap.String("request_protocol", ctx.Request.Proto),
		zap.Any("ctx_value_ajax", ctx.Value("ajax")),
	)
}
