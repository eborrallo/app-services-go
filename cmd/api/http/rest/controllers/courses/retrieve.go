package courses

import (
	"app-services-go/internal/courses/application/fetching"
	course "app-services-go/internal/courses/domain"
	"app-services-go/kit/query"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type retrieveRequest struct {
	ID string `json:"id" binding:"required" uri:"id"`
}

// RetrieveHandler returns an HTTP controllers for courses fetch.
func RetrieveHandler(queryBus query.Bus) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("RetrieveHandler ctx")
		var req retrieveRequest

		if err := ctx.BindUri(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		entity, err := queryBus.Ask(ctx, fetching.NewCourseQuery(req.ID))

		if err != nil {
			switch {
			case errors.Is(err, course.ErrInvalidCourseID),
				errors.Is(err, course.ErrEmptyCourseName), errors.Is(err, course.ErrInvalidCourseID):
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			default:
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}
		ctx.JSON(http.StatusOK, entity)

	}
}
