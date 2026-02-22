package handlers

import "go.uber.org/zap"

type MessageHandler struct {
	logger *zap.Logger
}

func NewMessageHandler(logger *zap.Logger) *MessageHandler {
	return &MessageHandler{
		logger: logger,
	}
}
