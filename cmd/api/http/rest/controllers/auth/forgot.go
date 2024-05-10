package auth

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
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := queryBus.Ask(ctx, fetching.NewUserByEmailQuery(req.Email))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = commandBus.Dispatch(ctx, forgot.NewForgotCommand(user.(domain.User).ID, user.(domain.User).Email))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.Status(http.StatusOK)
	}
}
