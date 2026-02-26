package types

import "time"

type RepeatType string

const (
	RepeatTypeOnce    RepeatType = "once"
	RepeatTypeDaily   RepeatType = "daily"
	RepeatTypeWeekly  RepeatType = "weekly"
	RepeatTypeMonthly RepeatType = "monthly"
	RepeatTypeYearly  RepeatType = "yearly"
)

type ScheduleMessage struct {
	Topic  string     `json:"topic"`
	Body   any        `json:"body"`
	Time   time.Time  `json:"time"`
	Repeat RepeatType `json:"repeat"`
}

func NewScheduleMessage(topic string, body any, time time.Time, repeat string) *ScheduleMessage {
	if topic == "" {
		return nil
	}
	// check if repeat is valid
	if repeat != string(RepeatTypeOnce) &&
		repeat != string(RepeatTypeDaily) &&
		repeat != string(RepeatTypeWeekly) &&
		repeat != string(RepeatTypeMonthly) &&
		repeat != string(RepeatTypeYearly) {
		return nil
	}
	return &ScheduleMessage{
		Topic:  topic,
		Body:   body,
		Time:   time,
		Repeat: RepeatType(repeat),
	}
}
