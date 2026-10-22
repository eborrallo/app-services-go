package auth_web2

import (
	"app-services-go/internal/auth/application/remember"
	"app-services-go/kit/command"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

type rememberRequest struct {
	NewPassword string `json:"new_password" binding:"required"`
}

func RememberHandler(commandBus command.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req rememberRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

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
		// Dispatch the DECODED token; the remember handler decrypts it.
		if err := commandBus.Dispatch(ctx.Request.Context(), remember.NewPasswordUserCommand(string(token), req.NewPassword)); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "could not update password"})
			return
		}
		ctx.Status(http.StatusOK)
	}
}
