package observability

import "github.com/gin-gonic/gin"

func ContextTraceMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.Next()
	}
}
