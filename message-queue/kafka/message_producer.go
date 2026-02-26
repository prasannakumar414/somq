package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/prasannakumar414/somq/types"
	"go.uber.org/zap"
)

type MessageProducer struct {
	logger *zap.Logger
	client sarama.SyncProducer
}

func NewMessageProducer(logger *zap.Logger, client sarama.SyncProducer) *MessageProducer {
	return &MessageProducer{
		logger: logger,
		client: client,
	}
}

// PublishMessage schedules a message to be published to Kafka at the time
// specified in message.Time. It returns immediately; the delay and send
// are handled in a background goroutine. If the scheduled time is in the
// past the message is published without any delay.
func (p *MessageProducer) PublishMessage(ctx context.Context, message *types.ScheduleMessage) {
	go func() {
		delay := time.Until(message.Time)
		if delay > 0 {
			p.logger.Info("Scheduling message",
				zap.String("topic", message.Topic),
				zap.Duration("delay", delay),
				zap.Time("scheduled_at", message.Time),
			)
			select {
			case <-time.After(delay):
				// Delay elapsed – proceed to publish.
			case <-ctx.Done():
				p.logger.Warn("Publish cancelled before scheduled time",
					zap.String("topic", message.Topic),
					zap.Error(ctx.Err()),
				)
				return
			}
		}

		body, err := json.Marshal(message.Body)
		if err != nil {
			p.logger.Error("Failed to marshal message body",
				zap.String("topic", message.Topic),
				zap.Error(err),
			)
			return
		}

		kafkaMsg := &sarama.ProducerMessage{
			Topic: message.Topic,
			Value: sarama.ByteEncoder(body),
		}

		partition, offset, err := p.client.SendMessage(kafkaMsg)
		if err != nil {
			p.logger.Error("Failed to publish message",
				zap.String("topic", message.Topic),
				zap.Error(err),
			)
			return
		}

		p.logger.Info("Message published",
			zap.String("topic", message.Topic),
			zap.Int32("partition", partition),
			zap.Int64("offset", offset),
		)

	}()
}
