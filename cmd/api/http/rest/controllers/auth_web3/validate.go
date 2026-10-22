package auth_web3

import (
	"app-services-go/internal/auth/application/creating"
	"app-services-go/internal/auth/application/fetching"
	"app-services-go/internal/auth/application/login"
	"app-services-go/internal/auth/domain"
	"app-services-go/kit/blokchain"
	"app-services-go/kit/command"
	"app-services-go/kit/crypt"
	"app-services-go/kit/query"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type validateMessaeRequest struct {
	Signature string `json:"signature" binding:"required"`
}

func ValidateMessageHandler(commandBus command.Bus, queryBus query.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		address := ctx.Param("address")

		var req validateMessaeRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if !blokchain.IsValidEthereumAddress(address) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Ethereum address"})
			return
		}

		result, err := queryBus.Ask(ctx.Request.Context(), fetching.NewUserBySignatureQuery(req.Signature, address))
		if err != nil {
			// A verified signature with no matching user means this is a new
			// web3 user: create it. Any other error is a failed verification.
			if errors.Is(err, domain.ErrUserNotFound) {
				token, refreshToken, tErr := issueTokens(domain.TokenPayload{ID: address})
				if tErr != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
					return
				}
				if err := commandBus.Dispatch(ctx.Request.Context(), creating.NewWeb3UserCommand(address, token, refreshToken)); err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
					return
				}
				ctx.JSON(http.StatusOK, gin.H{
					"token":         token,
					"refresh_token": refreshToken,
				})
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			return
		}

		user, ok := result.(domain.User)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected response"})
			return
		}

		token, refreshToken, tErr := issueTokens(domain.TokenPayload{ID: user.ID, Email: user.Email})
		if tErr != nil {
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

// issueTokens mints the access and refresh JWTs for the given payload.
func issueTokens(payload domain.TokenPayload) (string, string, error) {
	token, err := crypt.CreateToken(payload, 24*time.Hour)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := crypt.CreateToken(payload, 30*24*time.Hour)
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}
