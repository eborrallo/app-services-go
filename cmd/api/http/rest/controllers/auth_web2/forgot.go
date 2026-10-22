package auth_web2

import (
	"app-services-go/internal/auth/application/fetching"
	"app-services-go/internal/auth/application/forgot"
	"app-services-go/internal/auth/domain"
	"app-services-go/kit/command"
	"app-services-go/kit/query"
	"net/http"

	"github.com/gin-gonic/gin"
)

type forgotRequest struct {
	Email string `json:"email" binding:"required"`
}

func ForgotHandler(commandBus command.Bus, queryBus query.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req forgotRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		result, err := queryBus.Ask(ctx.Request.Context(), fetching.NewUserByEmailQuery(req.Email))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "could not process request"})
			return
		}

		user, ok := result.(domain.User)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "could not process request"})
			return
		}

		if err := commandBus.Dispatch(ctx.Request.Context(), forgot.NewForgotCommand(user.ID, user.Email)); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "could not process request"})
			return
		}

		ctx.Status(http.StatusOK)
	}
}
