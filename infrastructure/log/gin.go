package log

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// GinLogger
/**
 * @Description:
 * @return gin.HandlerFunc
 */
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		rid := c.GetHeader("X-Request-ID")

		// 用户链路调试
		if debugID := c.GetHeader(string(DebugRequestId)); debugID != "" {
			rid = debugID
		}

		if rid == "" {
			rid = uuid.NewString()
		}

		var ridCtx = context.WithValue(c.Request.Context(), RequestIdKey, rid) // session id
		c.Request = c.Request.WithContext(ridCtx)
		c.Next()
		cost := time.Since(start)
		GetLogger().Info(path,
			zap.String(SessionId.ToString(), rid),
			zap.Int(StatusCode.ToString(), c.Writer.Status()),
			zap.String(HttpMethod.ToString(), c.Request.Method),
			zap.String(HttpPath.ToString(), path),
			zap.String(QueryText.ToString(), query),
			zap.String(ClientIp.ToString(), c.ClientIP()),
			zap.String(UserAgent.ToString(), c.Request.UserAgent()),
			zap.String(Errors.ToString(), c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration(Duration.ToString(), cost),
		)
	}
}

// GinRecovery
/**
 * @Description:
 * @param stack
 * @return gin.HandlerFunc
 */
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				fmt.Println(err)
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if ok := errors.Is(ne.Err, se); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					GetLogger().Sugar().Error(c.Request.URL.Path, zap.Any(Errors.ToString(), err), zap.String(HttpRequest.ToString(), string(httpRequest)))
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: error check
					c.Abort()
					return
				}

				if stack {
					GetLogger().Sugar().Error("Recovery from panic", zap.Any(Errors.ToString(), err), zap.String(HttpRequest.ToString(), string(httpRequest)), zap.String("stack", string(debug.Stack())))
				} else {
					GetLogger().Sugar().Error("Recovery from panic", zap.Any(Errors.ToString(), err), zap.String(HttpRequest.ToString(), string(httpRequest)))
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
