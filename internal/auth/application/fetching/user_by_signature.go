package fetching

import (
	user "app-services-go/internal/auth/domain"
	"app-services-go/kit/blokchain"
	"app-services-go/kit/query"
	"context"
	"errors"
)

const UserBySignatureQueryType query.Type = "query.fetching.user_by_signature"

// UserBySignatureQuery is the query dispatched to create a new user.
type UserBySignatureQuery struct {
	signature string
	signer    string
}

// NewUserBySignatureQuery creates a new UserBySignatureQuery.
func NewUserBySignatureQuery(signature string, signer string) UserBySignatureQuery {
	return UserBySignatureQuery{
		signature: signature,
		signer:    signer,
	}
}

func (c UserBySignatureQuery) Type() query.Type {
	return UserBySignatureQueryType
}

// UserBySignatureQueryHandler is the query controllers responsible for creating users.
type UserBySignatureQueryHandler struct {
	userRepository        user.UserRepository
	userMessageRepository user.UserMessageRepository
	signatureVerificator  blokchain.SignatureVerificator
}

// NewUserBySignatureQueryHandler initializes a new UserBySignatureQueryHandler.
func NewUserBySignatureQueryHandler(userRepository user.UserRepository, userMessageRepository user.UserMessageRepository, signatureVerificator blokchain.SignatureVerificator) UserBySignatureQueryHandler {
	return UserBySignatureQueryHandler{
		userRepository:        userRepository,
		userMessageRepository: userMessageRepository,
		signatureVerificator:  signatureVerificator,
	}
}

// Handle implements the query.Handler interface.
func (h UserBySignatureQueryHandler) Handle(ctx context.Context, qry query.Query) (interface{}, error) {
	userQry, ok := qry.(UserBySignatureQuery)
	if !ok {
		return nil, errors.New("unexpected query")
	}
	message, err := h.userMessageRepository.GetMessage(userQry.signer)
	if err != nil {
		return nil, err
	}
	err = h.signatureVerificator.VerifyMessage(message, userQry.signature, userQry.signer)
	if err != nil {
		return nil, err
	}
	user, err := h.userRepository.FetchByAddress(ctx, userQry.signature)
	if err != nil {
		return nil, err
	}

	return user, nil
}
