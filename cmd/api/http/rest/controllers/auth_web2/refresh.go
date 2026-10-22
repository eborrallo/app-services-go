package auth_web2

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
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		var payload domain.TokenPayload
		if err := crypt.GetPayloadFromToken(req.RefreshToken, &payload); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		result, err := queryBus.Ask(ctx.Request.Context(), fetching.NewUserByIdQuery(payload.ID))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		user, ok := result.(domain.User)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		newPayload := domain.TokenPayload{ID: user.ID, Email: user.Email}
		token, err := crypt.CreateToken(newPayload, 24*time.Hour)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
			return
		}
		refreshToken, err := crypt.CreateToken(newPayload, 30*24*time.Hour)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
			return
		}

		if err := commandBus.Dispatch(ctx.Request.Context(), login.NewLoginCommand(token, refreshToken)); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not complete refresh"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}
