package auth_web3

import (
	"app-services-go/internal/auth/application/creating"
	"app-services-go/internal/auth/application/fetching"
	"app-services-go/internal/auth/application/login"
	"app-services-go/kit/blokchain"
	"app-services-go/kit/command"
	"app-services-go/kit/crypt"
	"app-services-go/kit/query"
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
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !blokchain.IsValidEthereumAddress(address) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Ethereum address"})
			return
		}

		user, err := queryBus.Ask(ctx, fetching.NewUserBySignatureQuery(req.Signature, address))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token := crypt.CreateToken(user, 24*time.Hour)
		refreshToken := crypt.CreateToken(user, 30*24*time.Hour)
		if user == nil {
			err := commandBus.Dispatch(ctx, creating.NewWeb3UserCommand(address, token, refreshToken))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}
		} else {
			err = commandBus.Dispatch(ctx, login.NewLoginCommand(token, refreshToken))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}
