package auth_web3

import (
	"app-services-go/internal/auth/application/message_generator"
	"app-services-go/kit/blokchain"
	"app-services-go/kit/command"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MessageHandler(commandBus command.Bus, messageGeneratorService message_generator.MessageGeneratorService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")
		address := ctx.Param("address")

		if !blokchain.IsValidEthereumAddress(address) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Ethereum address"})
			return
		}
		message, err := messageGeneratorService.Generate(address, origin)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		commandBus.Dispatch(ctx, message_generator.NewSaveMessageCommand(address, message))
		ctx.JSON(http.StatusOK, gin.H{
			"message": message,
		})
	}
}
