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
		if debugID := c.GetHeader(DEBUG_REQUEST_ID); debugID != "" {
			rid = debugID
		}

		if rid == "" {
			rid = uuid.NewString()
		}

		var ridCtx = context.WithValue(c.Request.Context(), REQUEST_ID_KEY, rid) // session id
		c.Request = c.Request.WithContext(ridCtx)
		c.Next()
		cost := time.Since(start)
		GetLogger().Info(path,
			zap.String(K_SessionId, rid),
			zap.Int(K_StatusCode, c.Writer.Status()),
			zap.String(K_HttpMethod, c.Request.Method),
			zap.String(K_HttpPath, path),
			zap.String(K_Query, query),
			zap.String(K_ClientIp, c.ClientIP()),
			zap.String(K_UserAgent, c.Request.UserAgent()),
			zap.String(K_Errors, c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration(K_Duration, cost),
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
					GetLogger().Sugar().Error(c.Request.URL.Path, zap.Any(K_Errors, err), zap.String(K_HttpRequest, string(httpRequest)))
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: error check
					c.Abort()
					return
				}

				if stack {
					GetLogger().Sugar().Error("Recovery from panic", zap.Any(K_Errors, err), zap.String(K_HttpRequest, string(httpRequest)), zap.String("stack", string(debug.Stack())))
				} else {
					GetLogger().Sugar().Error("Recovery from panic", zap.Any(K_Errors, err), zap.String(K_HttpRequest, string(httpRequest)))
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
