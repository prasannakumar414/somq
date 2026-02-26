package services

import (
	"context"
	"errors"
	"time"

	"github.com/prasannakumar414/somq/types"
	"go.uber.org/zap"
)

type MessageService struct {
	logger          *zap.Logger
	messageDAO      ScheduleMessageDAO
	messageProducer MessageProducer
}

func NewMessageService(logger *zap.Logger, messageDAO ScheduleMessageDAO, messageProducer MessageProducer) *MessageService {
	return &MessageService{
		logger:          logger,
		messageDAO:      messageDAO,
		messageProducer: messageProducer,
	}
}

func (s *MessageService) ScheduleMessage(ctx context.Context, message *types.ScheduleMessage) error {
	if message.Time.Before(time.Now()) {
		return errors.New("scheduled time is in the past")
	}
	if message.Time.After(time.Now().Add(24 * time.Hour)) {
		return s.messageDAO.CreateScheduleMessage(ctx, message)
	}
	s.messageProducer.PublishMessage(ctx, message)
	return nil
}
