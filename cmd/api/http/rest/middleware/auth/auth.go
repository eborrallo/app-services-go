package auth

import (
	"app-services-go/kit/crypt"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UnsignedResponse struct {
	Message interface{} `json:"message"`
}

// Middleware creates a Gin middleware for JWT verification
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
				Message: err.Error(),
			})
			return
		}
		err = crypt.VerifyToken(jwtToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
				Message: "bad jwt token",
			})
			return
		}

		c.Next()
	}
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}
