package auth_web3

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"app-services-go/internal/auth/domain"
	"app-services-go/kit/command/commandmocks"
	"app-services-go/kit/query/querymocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_ValdiateUser(t *testing.T) {
	commandBus := new(commandmocks.Bus)
	queryBus := new(querymocks.Bus)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/validate/:address", ValidateMessageHandler(commandBus, queryBus))

	t.Run("given an non signature it returns 400", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodPost, "/validate/123", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
	t.Run("given an invalid address it returns 400", func(t *testing.T) {
		validateReq := validateMessaeRequest{
			Signature: "123",
		}

		b, err := json.Marshal(validateReq)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/validate/123", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
	t.Run("given an invalid signature it returns 400", func(t *testing.T) {
		queryBus.On(
			"Ask",
			mock.Anything,
			mock.AnythingOfType("fetching.UserBySignatureQuery"),
		).Return(nil, errors.New("something unexpected happened")).Once()

		validateReq := validateMessaeRequest{
			Signature: "123",
		}

		b, err := json.Marshal(validateReq)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/validate/0x2bbd10fc8793c35f25de6f25ed1d6cff1402b473", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
	t.Run("given an existing user update tokens and returns 200", func(t *testing.T) {
		user := domain.User{
			ID:       "8a1c5cdc-ba57-445a-994d-aa412d23723f",
			Email:    "aaa@gmail.com",
			Password: "123",
		}
		queryBus.On(
			"Ask",
			mock.Anything,
			mock.AnythingOfType("fetching.UserBySignatureQuery"),
		).Return(user, nil).Once()

		commandBus.On(
			"Dispatch",
			mock.Anything,
			mock.AnythingOfType("login.LoginCommand"),
		).Return(nil).Once()

		validateReq := validateMessaeRequest{
			Signature: "123",
		}

		b, err := json.Marshal(validateReq)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/validate/0x2bbd10fc8793c35f25de6f25ed1d6cff1402b473", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("given an non existing user update tokens and returns 200", func(t *testing.T) {
		queryBus.On(
			"Ask",
			mock.Anything,
			mock.AnythingOfType("fetching.UserBySignatureQuery"),
		).Return(nil, nil).Once()

		commandBus.On(
			"Dispatch",
			mock.Anything,
			mock.AnythingOfType("creating.Web3UserCommand"),
		).Return(nil).Once()

		validateReq := validateMessaeRequest{
			Signature: "123",
		}

		b, err := json.Marshal(validateReq)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/validate/0x2bbd10fc8793c35f25de6f25ed1d6cff1402b473", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
