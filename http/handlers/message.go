package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/prasannakumar414/somq/services"
	"github.com/prasannakumar414/somq/types"
	"go.uber.org/zap"
)

type MessageHandler struct {
	logger         *zap.Logger
	messageService *services.MessageService
}

func NewMessageHandler(logger *zap.Logger, messageService *services.MessageService) *MessageHandler {
	return &MessageHandler{
		logger:         logger,
		messageService: messageService,
	}
}

func (h *MessageHandler) ScheduleMessage(w http.ResponseWriter, r *http.Request) (error, int, any) {
	var message types.ScheduleMessage
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		return err, 0, nil
	}
	err = h.messageService.ScheduleMessage(r.Context(), &message)
	if err != nil {
		return err, 0, nil
	}
	return nil, 200, "Message scheduled successfully"
}
