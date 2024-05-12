package fetching

import (
	user "app-services-go/internal/auth/domain"
	"app-services-go/kit/query"
	"context"
	"errors"
)

const UserByEmailQueryType query.Type = "query.fetching.user_by_email"

// UserByEmailQuery is the query dispatched to create a new user.
type UserByEmailQuery struct {
	email string
}

// NewUserByEmailQuery creates a new UserByEmailQuery.
func NewUserByEmailQuery(email string) UserByEmailQuery {
	return UserByEmailQuery{
		email: email,
	}
}

func (c UserByEmailQuery) Type() query.Type {
	return UserByEmailQueryType
}

// UserByEmailQueryHandler is the query controllers responsible for creating users.
type UserByEmailQueryHandler struct {
	userRepository user.UserRepository
}

// NewUserByEmailQueryHandler initializes a new UserByEmailQueryHandler.
func NewUserByEmailQueryHandler(userRepository user.UserRepository) UserByEmailQueryHandler {
	return UserByEmailQueryHandler{
		userRepository: userRepository,
	}
}

// Handle implements the query.Handler interface.
func (h UserByEmailQueryHandler) Handle(ctx context.Context, qry query.Query) (interface{}, error) {
	userQry, ok := qry.(UserByEmailQuery)
	if !ok {
		return nil, errors.New("unexpected query")
	}

	user, err := h.userRepository.FetchByEmail(ctx, userQry.email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
