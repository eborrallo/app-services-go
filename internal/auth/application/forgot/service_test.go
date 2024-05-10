package forgot

import (
	"app-services-go/internal/auth/infrastructure/email/smptmocks"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ForgotUserService_Send_Succeed(t *testing.T) {

	senderMock := new(smptmocks.Sender)
	userId := "37a0f027-15e6-47cc-a5d2-64183281087e"
	email := "aaa@gmail.com"

	senderMock.On("Send", email, mock.AnythingOfType("string")).Return(nil)

	service := NewForgotUserService(senderMock)
	err := service.SendForgotEmail(context.Background(), userId, email)
	assert.NoError(t, err)

	senderMock.AssertCalled(t, "Send", email, mock.AnythingOfType("string"))
}
