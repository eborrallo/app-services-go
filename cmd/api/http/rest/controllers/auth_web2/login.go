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

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginHandler(commandBus command.Bus, queryBus query.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req loginRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		result, err := queryBus.Ask(ctx.Request.Context(), fetching.NewUserByEmailQuery(req.Email))
		if err != nil {
			// Do not reveal whether the email exists: same 401 as a bad password.
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		user, ok := result.(domain.User)
		if !ok || !crypt.Compare(req.Password, user.Password) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		payload := domain.TokenPayload{ID: user.ID, Email: user.Email}
		token, err := crypt.CreateToken(payload, 24*time.Hour)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
			return
		}
		refreshToken, err := crypt.CreateToken(payload, 30*24*time.Hour)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
			return
		}

		if err := commandBus.Dispatch(ctx.Request.Context(), login.NewLoginCommand(token, refreshToken)); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not complete login"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}
