package mongo

import (
	"context"
	"time"

	"github.com/prasannakumar414/somq/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type ScheduleMessageDAO struct {
	client *mongo.Client
	logger *zap.Logger
}

func NewScheduleMessageDAO(client *mongo.Client, logger *zap.Logger) *ScheduleMessageDAO {
	return &ScheduleMessageDAO{
		client: client,
		logger: logger,
	}
}

// CreateScheduleMessage creates a new schedule message.
// Stores in a collection based on the repeat type:
//   - once     → schedule_messages_once
//   - daily    → schedule_messages_daily
//   - weekly   → schedule_messages_weekly
//   - monthly  → schedule_messages_monthly
//   - yearly   → schedule_messages_yearly
func (s *ScheduleMessageDAO) CreateScheduleMessage(ctx context.Context, scheduleMessage *types.ScheduleMessage) error {
	collectionName := collectionForRepeatType(scheduleMessage.Repeat)
	collection := s.client.Database("somq").Collection(collectionName)
	_, err := collection.InsertOne(ctx, scheduleMessage)
	if err != nil {
		s.logger.Error("Failed to create schedule message", zap.Error(err), zap.String("collection", collectionName))
		return err
	}
	return nil
}

// Get Messages Scheduled today
func (s *ScheduleMessageDAO) GetMessagesScheduledToday(ctx context.Context) ([]types.ScheduleMessage, error) {
	collectionName := collectionForRepeatType(types.RepeatTypeOnce)
	collection := s.client.Database("somq").Collection(collectionName)
	var scheduleMessages []types.ScheduleMessage
	today := time.Now().Truncate(24 * time.Hour)
	cursor, err := collection.Find(ctx, bson.M{"time": bson.M{"$gte": today, "$lt": today.Add(24 * time.Hour)}})
	if err != nil {
		s.logger.Error("Failed to get schedule messages", zap.Error(err), zap.String("collection", collectionName))
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &scheduleMessages); err != nil {
		s.logger.Error("Failed to decode schedule messages", zap.Error(err), zap.String("collection", collectionName))
		return nil, err
	}
	return scheduleMessages, nil
}

// Delete Schedule Message
func (s *ScheduleMessageDAO) DeleteScheduleMessage(ctx context.Context, id string, repeat types.RepeatType) error {
	collectionName := collectionForRepeatType(repeat)
	collection := s.client.Database("somq").Collection(collectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		s.logger.Error("Failed to delete schedule message", zap.Error(err), zap.String("collection", collectionName))
		return err
	}
	return nil
}

// collectionForRepeatType returns the MongoDB collection name for the given repeat type.
func collectionForRepeatType(repeat types.RepeatType) string {
	switch repeat {
	case types.RepeatTypeOnce:
		return "schedule_messages_once"
	case types.RepeatTypeDaily:
		return "schedule_messages_daily"
	case types.RepeatTypeWeekly:
		return "schedule_messages_weekly"
	case types.RepeatTypeMonthly:
		return "schedule_messages_monthly"
	case types.RepeatTypeYearly:
		return "schedule_messages_yearly"
	default:
		return "schedule_messages_once"
	}
}
