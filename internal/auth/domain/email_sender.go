package domain

//go:generate mockery --case=snake --outpkg=smptmocks --output=../infrastructure/email/smptmocks --name=Sender

type Sender interface {
	Send(email string, content string)
}
