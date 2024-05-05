package creating

import (
	user "app-services-go/internal/auth/domain"
	"app-services-go/internal/auth/infrastructure/storage/storagemocks"
	"app-services-go/kit/event"
	"app-services-go/kit/event/eventmocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_UserService_CreateUser_RepositoryError(t *testing.T) {
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	userName := "Test User"
	userEmail := "aaa@gmail.com"
	userPassword := "123123"

	userRepositoryMock := new(storagemocks.UserRepository)
	userRepositoryMock.On("Save", mock.Anything, mock.AnythingOfType("User")).Return(errors.New("something unexpected happened"))

	eventBusMock := new(eventmocks.Bus)

	userService := NewUserService(userRepositoryMock, eventBusMock)

	err := userService.CreateUser(context.Background(), userID, userName, userEmail, userPassword)
	userRepositoryMock.AssertExpectations(t)
	eventBusMock.AssertExpectations(t)
	assert.Error(t, err)
}

func Test_UserService_CreateUser_EventsBusError(t *testing.T) {
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	userName := "Test User"
	userEmail := "aaa@gmail.com"
	userPassword := "123123"

	userRepositoryMock := new(storagemocks.UserRepository)
	userRepositoryMock.On("Save", mock.Anything, mock.AnythingOfType("User")).Return(nil)

	eventBusMock := new(eventmocks.Bus)
	eventBusMock.On("Publish", mock.Anything, mock.AnythingOfType("[]event.Event")).Return(errors.New("something unexpected happened"))

	userService := NewUserService(userRepositoryMock, eventBusMock)

	err := userService.CreateUser(context.Background(), userID, userName, userEmail, userPassword)

	userRepositoryMock.AssertExpectations(t)
	eventBusMock.AssertExpectations(t)
	assert.Error(t, err)
}

func Test_UserService_CreateUser_Succeed(t *testing.T) {
	userID := "37a0f027-15e6-47cc-a5d2-64183281087e"
	userName := "Test User"
	userEmail := "aaa@gmail.com"
	userPassword := "123123"

	userRepositoryMock := new(storagemocks.UserRepository)
	userRepositoryMock.On("Save", mock.Anything, mock.AnythingOfType("User")).Return(nil)

	eventBusMock := new(eventmocks.Bus)
	eventBusMock.On("Publish", mock.Anything, mock.MatchedBy(func(events []event.Event) bool {
		evt := events[0].(user.UserCreatedEvent)
		return evt.UserName() == userName
	})).Return(nil)

	eventBusMock.On("Publish", mock.Anything, mock.AnythingOfType("[]event.Event")).Return(nil)

	userService := NewUserService(userRepositoryMock, eventBusMock)

	err := userService.CreateUser(context.Background(), userID, userName, userEmail, userPassword)

	userRepositoryMock.AssertExpectations(t)
	eventBusMock.AssertExpectations(t)
	assert.NoError(t, err)
}
