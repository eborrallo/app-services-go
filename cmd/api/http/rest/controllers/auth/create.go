package auth

import (
	"app-services-go/internal/auth/application/creating"
	"app-services-go/kit/command"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createRequest struct {
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" validate:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateHandler(commandBus command.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req createRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		err := commandBus.Dispatch(ctx, creating.NewUserCommand(
			req.ID,
			req.Name,
			req.Email,
			req.Password,
		))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		ctx.Status(http.StatusCreated)
	}
}
