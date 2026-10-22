package auth_web2

import (
	"app-services-go/internal/auth/application/validating"
	"app-services-go/kit/command"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateHandler(commandBus command.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		base64Token := ctx.Param("token")
		token, err := base64.URLEncoding.DecodeString(base64Token)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}
		if string(token) == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}
		// Dispatch the DECODED token; the validating handler decrypts it.
		if err := commandBus.Dispatch(ctx.Request.Context(), validating.NewValidateUserCommand(string(token))); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "could not validate user"})
			return
		}
		ctx.Status(http.StatusOK)
	}
}
