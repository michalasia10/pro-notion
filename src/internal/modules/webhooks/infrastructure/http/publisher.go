package http

import (
	"encoding/json"
	"log"

	"github.com/ThreeDotsLabs/watermill/message"

	shared "src/internal/modules/shared/domain"
	sharedEvents "src/internal/modules/shared/domain/events"
	"src/internal/modules/webhooks/domain"
)

// WatermillEventPublisher implements WebhookEventPublisher using Watermill
type WatermillEventPublisher struct {
	publisher message.Publisher
	logger    *log.Logger
}

// NewWatermillEventPublisher creates a new Watermill event publisher
func NewWatermillEventPublisher(publisher message.Publisher, logger *log.Logger) *WatermillEventPublisher {
	return &WatermillEventPublisher{
		publisher: publisher,
		logger:    logger,
	}
}

// PublishEvent publishes a webhook event to Watermill
func (p *WatermillEventPublisher) PublishEvent(event *domain.WebhookEvent) error {
	// Convert domain event to shared event
	sharedEvent := sharedEvents.NotionWebhookReceived{
		Payload: []byte(event.Payload),
	}

	eventBytes, err := json.Marshal(sharedEvent)
	if err != nil {
		p.logger.Printf("Failed to marshal webhook event: %v", err)
		return domain.ErrProcessingFailed
	}

	msg := message.NewMessage(shared.NewUUIDGenerator().NewID("webhook"), eventBytes)

	if err := p.publisher.Publish(sharedEvents.NotionWebhookReceivedTopic, msg); err != nil {
		p.logger.Printf("Failed to publish webhook event to Watermill: %v", err)
		return domain.WebhookProcessingError{
			Code:    "PUBLISH_FAILED",
			Message: err.Error(),
		}
	}

	p.logger.Printf("Successfully published webhook event to topic: %s", sharedEvents.NotionWebhookReceivedTopic)
	return nil
}
