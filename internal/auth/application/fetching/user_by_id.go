package fetching

import (
	user "app-services-go/internal/auth/domain"
	"app-services-go/kit/query"
	"context"
	"errors"
)

const UserByIdQueryType query.Type = "query.fetching.user_by_id"

// UserByIdQuery is the query dispatched to create a new user.
type UserByIdQuery struct {
	email string
}

// NewUserByIdQuery creates a new UserByIdQuery.
func NewUserByIdQuery(email string) UserByIdQuery {
	return UserByIdQuery{
		email: email,
	}
}

func (c UserByIdQuery) Type() query.Type {
	return UserByIdQueryType
}

// UserByIdQueryHandler is the query controllers responsible for creating users.
type UserByIdQueryHandler struct {
	userRepository user.UserRepository
}

// NewUserByIdQueryHandler initializes a new UserByIdQueryHandler.
func NewUserByIdQueryHandler(userRepository user.UserRepository) UserByIdQueryHandler {
	return UserByIdQueryHandler{
		userRepository: userRepository,
	}
}

// Handle implements the query.Handler interface.
func (h UserByIdQueryHandler) Handle(ctx context.Context, qry query.Query) (interface{}, error) {
	userQry, ok := qry.(UserByIdQuery)
	if !ok {
		return nil, errors.New("unexpected query")
	}

	user, err := h.userRepository.FetchById(ctx, userQry.email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
