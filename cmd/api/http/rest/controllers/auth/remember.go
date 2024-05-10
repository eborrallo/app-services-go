package auth

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
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		base64Token := ctx.Param("token")
		token, err := base64.URLEncoding.DecodeString(base64Token)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "invalid token")
			return
		}
		if string(token) == "" {
			ctx.JSON(http.StatusBadRequest, "invalid token")
			return
		}
		err = commandBus.Dispatch(ctx, remember.NewPasswordUserCommand(string(base64Token), req.NewPassword))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		ctx.Status(http.StatusOK)
	}
}
