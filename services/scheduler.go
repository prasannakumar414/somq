package services

import (
	"context"
	"time"

	"github.com/prasannakumar414/somq/types"
	"go.uber.org/zap"
)

type SchedulerService struct {
	scheduleMessageDAO ScheduleMessageDAO
	logger             *zap.Logger
	messageProducer    MessageProducer
}

// an interface for dao
type ScheduleMessageDAO interface {
	CreateScheduleMessage(ctx context.Context, message *types.ScheduleMessage) error
	GetMessagesScheduledToday(ctx context.Context) ([]types.ScheduleMessage, error)
	DeleteScheduleMessage(ctx context.Context, id string, repeat types.RepeatType) error
}

// an interface for message producer
type MessageProducer interface {
	PublishMessage(ctx context.Context, message *types.ScheduleMessage)
}

func NewSchedulerService(scheduleMessageDAO ScheduleMessageDAO, logger *zap.Logger, messageProducer MessageProducer) *SchedulerService {
	return &SchedulerService{
		scheduleMessageDAO: scheduleMessageDAO,
		logger:             logger,
		messageProducer:    messageProducer,
	}
}

// Method that runs every day at midnight and publishes messages scheduled for today
func (s *SchedulerService) Run(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		messages, err := s.scheduleMessageDAO.GetMessagesScheduledToday(ctx)
		if err != nil {
			s.logger.Error("Failed to get schedule messages", zap.Error(err))
			continue
		}
		for _, message := range messages {
			s.messageProducer.PublishMessage(ctx, &message)
		}
	}
}
