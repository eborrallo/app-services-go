package auth

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

func TestHandler_Login_ServiceError(t *testing.T) {
	commandBus := new(commandmocks.Bus)
	queryBus := new(querymocks.Bus)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", LoginHandler(commandBus, queryBus))

	t.Run("given an invalid request it returns 400", func(t *testing.T) {

		loginReq := createRequest{
			Password: "123",
		}

		b, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
	t.Run("given an invalid email it returns 400", func(t *testing.T) {
		queryBus.On(
			"Ask",
			mock.Anything,
			mock.AnythingOfType("fetching.UserByEmailQuery"),
		).Return(nil, errors.New("something unexpected happened")).Once()

		loginReq := createRequest{
			Email:    "aaa@aa.com",
			Password: "123",
		}

		b, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given an error new login it returns 400", func(t *testing.T) {
		queryBus.On(
			"Ask",
			mock.Anything,
			mock.AnythingOfType("fetching.UserByEmailQuery"),
		).Return(domain.User{}, nil).Once()
		commandBus.On(
			"Dispatch",
			mock.Anything,
			mock.AnythingOfType("login.LoginCommand"),
		).Return(errors.New("something unexpected happened")).Once()

		loginReq := createRequest{
			Email:    "aaa@aa.com",
			Password: "123",
		}

		b, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given a valid request it returns 201", func(t *testing.T) {
		commandBus.On(
			"Dispatch",
			mock.Anything,
			mock.AnythingOfType("login.LoginCommand"),
		).Return(nil).Once()
		queryBus.On(
			"Ask",
			mock.Anything,
			mock.AnythingOfType("fetching.UserByEmailQuery"),
		).Return(domain.User{}, nil).Once()

		loginReq := createRequest{
			Email:    "aaa@aa.com",
			Password: "123",
		}

		b, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(b))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

}
