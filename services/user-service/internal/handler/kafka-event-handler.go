package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-ecommerce-application/libs/kafka/events"
	"github.com/go-ecommerce-application/services/user-service/internal/service"
)

type KafkaEventHandler struct {
	userProfileService service.UserProfileService
}

func NewKafkaEventHandler(userProfileService service.UserProfileService) *KafkaEventHandler {
	return &KafkaEventHandler{
		userProfileService: userProfileService,
	}
}

func (h *KafkaEventHandler) HandleUserSignedUpEvent(ctx context.Context, messageValue []byte) error {
	var event events.UserSignedUp
	if err := json.Unmarshal(messageValue, &event); err != nil {
		log.Printf("error unmarshaling user signed up event: %v", err)
		return err
	}

	log.Printf("received user signed up event: user_id=%s, email=%s", event.UserID, event.Email)

	if err := h.userProfileService.HandleUserSignedUpEvent(&event); err != nil {
		log.Printf("error handling user signed up event: %v", err)
		return err
	}

	return nil
}
