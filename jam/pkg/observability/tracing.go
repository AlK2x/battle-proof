package observability

import (
	"log/slog"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
)

func ContextTraceMiddleware(appName string) gin.HandlerFunc {
	return otelgin.Middleware(appName)
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		spanCtx := trace.SpanContextFromContext(c.Request.Context())
		latency := time.Since(start)
		slog.InfoContext(c.Request.Context(),
			"Request Completed",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", latency),
			slog.String("trace_id", spanCtx.TraceID().String()),
		)
	}
}

type Singleton struct{}

var instance *Singleton
var once sync.Once

func GetInstance() *Singleton {
	if instance != nil {
		return instance
	}
	once.Do(func() {
		if instance != nil {
			return
		}
		instance = &Singleton{}
	})
	return instance
}
