package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"app-services-go/internal/auth/application/validating"
	"app-services-go/kit/command/commandmocks"
	"app-services-go/kit/crypt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_ValdiateUser(t *testing.T) {
	commandBus := new(commandmocks.Bus)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/validate/:token", ValidateHandler(commandBus))

	t.Run("given an invalid request it returns 400", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, "/validate/123", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given a valid request it returns 200", func(t *testing.T) {
		userId := "8a1c5cdc-ba57-445a-994d-aa412d23723f"
		encryptedId := crypt.Encrypt(userId)
		commandBus.On(
			"Dispatch",
			mock.Anything,
			mock.MatchedBy(func(cmd validating.ValidateUserCommand) bool {
				return cmd.Token == encryptedId
			})).Return(nil)

		req, err := http.NewRequest(http.MethodGet, "/validate/"+encryptedId, nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
