package message_generator

import (
	"app-services-go/configs"
	"app-services-go/internal/auth/domain"
	"app-services-go/kit/event"
	"time"
)

type MessageGeneratorService struct {
}

func NewMessageGeneratorService() MessageGeneratorService {
	return MessageGeneratorService{}
}

func (s MessageGeneratorService) Generate(address string, origin string) (string, error) {
	config, err := configs.GetBlokchainConfig()
	if err != nil {
		return "", err
	}
	message := "I wants you to sign in with your Ethereum account: " + address + " Sign in with ethereum "
	message = message + " Version: 1 Network: " + config.Network + " Chain: " + config.Chain + " ChainId: " + config.ChainId
	message = message + " Timestamp: " + time.Now().Local().String()
	return message, nil
}

type SaveMessageService struct {
	UserMessageRepository domain.UserMessageRepository
	eventBus              event.Bus
}

func NewSaveMessageService(userMessageRepository domain.UserMessageRepository, eventBus event.Bus) SaveMessageService {
	return SaveMessageService{
		UserMessageRepository: userMessageRepository,
		eventBus:              eventBus,
	}
}

func (s SaveMessageService) Save(address, message string) {

	s.UserMessageRepository.SaveMessage(address, message)
}
