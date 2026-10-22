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
			c.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{
				Message: err.Error(),
			})
			return
		}
		err = crypt.VerifyToken(jwtToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{
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

	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	// The scheme must be "Bearer"; otherwise parts[1] is not a bearer token.
	if parts[0] != "Bearer" {
		return "", errors.New("authorization header must use the Bearer scheme")
	}

	return parts[1], nil
}
