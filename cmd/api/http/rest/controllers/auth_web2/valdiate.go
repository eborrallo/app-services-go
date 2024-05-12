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
			ctx.JSON(http.StatusBadRequest, "invalid token")
			return
		}
		if string(token) == "" {
			ctx.JSON(http.StatusBadRequest, "invalid token")
			return
		}
		err = commandBus.Dispatch(ctx, validating.NewValidateUserCommand(string(base64Token)))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		ctx.Status(http.StatusOK)
	}
}
