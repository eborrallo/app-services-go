package validating

import (
	"app-services-go/internal/auth/domain"
	"app-services-go/internal/auth/domain/events"
	"app-services-go/internal/auth/infrastructure/storage/storagemocks"
	"app-services-go/kit/event"
	"app-services-go/kit/event/eventmocks"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ValidateUserService_Valdiate_Succeed(t *testing.T) {
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	userName := "Test User"
	userEmail := "aaa@gmail.com"
	userPassword := "123123"
	user := domain.User{
		ID:        userID,
		Name:      userName,
		Email:     userEmail,
		Password:  userPassword,
		Validated: false,
	}

	userRepositoryMock := new(storagemocks.UserRepository)
	userRepositoryMock.On("FetchById", mock.Anything, userID).Return(user, nil)
	userRepositoryMock.On("Update", mock.Anything, mock.MatchedBy(func(user domain.User) bool {
		return user.Validated == true
	})).Return(nil)

	eventBusMock := new(eventmocks.Bus)
	eventBusMock.On("Publish", mock.Anything, mock.MatchedBy(func(_events []event.Event) bool {
		evt := _events[0].(events.UserValidatedEvent)
		return evt.Email == userEmail && evt.UserName() == userName && evt.UserID() == userID
	})).Return(nil)

	eventBusMock.On("Publish", mock.Anything, mock.AnythingOfType("[]event.Event")).Return(nil)

	userService := NewValidateUserService(userRepositoryMock, eventBusMock)

	err := userService.ValidateUser(context.Background(), userID)
	userRepositoryMock.AssertExpectations(t)
	eventBusMock.AssertExpectations(t)
	assert.NoError(t, err)
}
