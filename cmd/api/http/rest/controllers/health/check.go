package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CheckHandler returns an HTTP controllers to perform health checks.
func CheckHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response := struct{ Status string }{Status: "ok"}
		ctx.JSON(http.StatusOK, response)
	}
}
