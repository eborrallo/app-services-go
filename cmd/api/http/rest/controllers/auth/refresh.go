package auth

import (
	"app-services-go/internal/auth/application/fetching"
	"app-services-go/internal/auth/application/login"
	"app-services-go/internal/auth/domain"
	"app-services-go/kit/command"
	"app-services-go/kit/crypt"
	"app-services-go/kit/query"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func RefreshHandler(commandBus command.Bus, queryBus query.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req refreshRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tokenUser := domain.User{}
		err := crypt.GetPayloadFromToken(req.RefreshToken, &tokenUser)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := queryBus.Ask(ctx, fetching.NewUserByEmailQuery(tokenUser.Email))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token := crypt.CreateToken(user, 24*time.Hour)
		refreshToken := crypt.CreateToken(user, 30*24*time.Hour)

		err = commandBus.Dispatch(ctx, login.NewLoginCommand(token, refreshToken))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}
