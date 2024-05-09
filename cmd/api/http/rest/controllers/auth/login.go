package auth

import (
	"app-services-go/internal/auth/application/fetching"
	"app-services-go/internal/auth/application/login"
	"app-services-go/kit/command"
	"app-services-go/kit/crypt"
	"app-services-go/kit/query"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateHandler returns an HTTP controllers for user auth creation.
func LoginHandler(commandBus command.Bus, queryBus query.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req loginRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := queryBus.Ask(ctx, fetching.NewUserByEmailQuery(req.Email))
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
