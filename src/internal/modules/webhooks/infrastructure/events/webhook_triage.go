package events

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"

	sharedEvents "src/internal/modules/shared/domain/events"
	sharedDTO "src/internal/modules/shared/interfaces/dto"
)

// NotionWebhookPayload represents the structure of a webhook payload from Notion
type NotionWebhookPayload struct {
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// WebhookTriage processes raw Notion webhook events and publishes specific domain events
type WebhookTriage struct {
	publisher message.Publisher
	logger    *log.Logger
}

// NewWebhookTriage creates a new webhook triage processor
func NewWebhookTriage(publisher message.Publisher, logger *log.Logger) *WebhookTriage {
	return &WebhookTriage{
		publisher: publisher,
		logger:    logger,
	}
}

// ProcessWebhook handles incoming webhook messages from Watermill
func (wt *WebhookTriage) ProcessWebhook(msg *message.Message) error {
	defer func() {
		msg.Ack()
	}()

	// Parse the raw webhook payload
	var webhookEvent sharedEvents.NotionWebhookReceived
	if err := json.Unmarshal(msg.Payload, &webhookEvent); err != nil {
		wt.logger.Printf("Failed to unmarshal webhook event: %v", err)
		return err
	}

	// Parse the Notion webhook payload
	var notionPayload NotionWebhookPayload
	if err := json.Unmarshal(webhookEvent.Payload, &notionPayload); err != nil {
		wt.logger.Printf("Failed to unmarshal Notion payload: %v", err)
		return err
	}

	// Process based on webhook type
	switch notionPayload.Type {
	case "page.updated":
		return wt.handlePageUpdated(notionPayload)
	case "page.created":
		return wt.handlePageCreated(notionPayload)
	case "page.deleted":
		return wt.handlePageDeleted(notionPayload)
	case "database.updated":
		return wt.handleDatabaseUpdated(notionPayload)
	default:
		wt.logger.Printf("Unknown webhook type: %s", notionPayload.Type)
		return nil // Ignore unknown types
	}
}

// handlePageUpdated processes page update events
func (wt *WebhookTriage) handlePageUpdated(payload NotionWebhookPayload) error {
	pageID, ok := payload.Data["id"].(string)
	if !ok {
		wt.logger.Printf("Missing page ID in page.updated event")
		return nil
	}

	databaseID, ok := payload.Data["database_id"].(string)
	if !ok {
		wt.logger.Printf("Missing database ID in page.updated event")
		return nil
	}

	// For now, assume all page updates in databases are task updates
	// TODO: Add more sophisticated logic to determine if it's a task vs other page types
	taskID := uuid.New() // Generate a task ID - in real implementation this would be looked up from database

	event := sharedEvents.NewTaskPropertiesUpdated(taskID, databaseID, pageID, payload.Data)

	return wt.publishTaskPropertiesUpdated(event)
}

// handlePageCreated processes page creation events
func (wt *WebhookTriage) handlePageCreated(payload NotionWebhookPayload) error {
	pageID, ok := payload.Data["id"].(string)
	if !ok {
		wt.logger.Printf("Missing page ID in page.created event")
		return nil
	}

	databaseID, ok := payload.Data["database_id"].(string)
	if !ok {
		wt.logger.Printf("Missing database ID in page.created event")
		return nil
	}

	// For now, assume all new pages in databases are tasks
	taskID := uuid.New()

	event := sharedEvents.NewTaskCreated(taskID, databaseID, pageID, payload.Data)

	return wt.publishTaskCreated(event)
}

// handlePageDeleted processes page deletion events
func (wt *WebhookTriage) handlePageDeleted(payload NotionWebhookPayload) error {
	pageID, ok := payload.Data["id"].(string)
	if !ok {
		wt.logger.Printf("Missing page ID in page.deleted event")
		return nil
	}

	databaseID, ok := payload.Data["database_id"].(string)
	if !ok {
		wt.logger.Printf("Missing database ID in page.deleted event")
		return nil
	}

	// For now, assume all deleted pages in databases are tasks
	taskID := uuid.New()

	event := sharedEvents.NewTaskDeleted(taskID, databaseID, pageID)

	return wt.publishTaskDeleted(event)
}

// handleDatabaseUpdated processes database update events
func (wt *WebhookTriage) handleDatabaseUpdated(payload NotionWebhookPayload) error {
	databaseID, ok := payload.Data["id"].(string)
	if !ok {
		wt.logger.Printf("Missing database ID in database.updated event")
		return nil
	}

	pageID, ok := payload.Data["page_id"].(string)
	if !ok {
		pageID = "" // Page ID might not be present for database updates
	}

	event := sharedEvents.NewDatabasePropertiesUpdated(databaseID, pageID, payload.Data)

	return wt.publishDatabasePropertiesUpdated(event)
}

// publishTaskPropertiesUpdated publishes TaskPropertiesUpdated event
func (wt *WebhookTriage) publishTaskPropertiesUpdated(event *sharedEvents.TaskPropertiesUpdated) error {
	dto := sharedDTO.TaskPropertiesUpdatedDTOFromDomain(event)
	return wt.publishDTO(sharedEvents.TaskPropertiesUpdatedTopic, dto)
}

// publishTaskCreated publishes TaskCreated event
func (wt *WebhookTriage) publishTaskCreated(event *sharedEvents.TaskCreated) error {
	dto := sharedDTO.TaskCreatedDTOFromDomain(event)
	return wt.publishDTO(sharedEvents.TaskCreatedTopic, dto)
}

// publishTaskDeleted publishes TaskDeleted event
func (wt *WebhookTriage) publishTaskDeleted(event *sharedEvents.TaskDeleted) error {
	dto := sharedDTO.TaskDeletedDTOFromDomain(event)
	return wt.publishDTO(sharedEvents.TaskDeletedTopic, dto)
}

// publishDatabasePropertiesUpdated publishes DatabasePropertiesUpdated event
func (wt *WebhookTriage) publishDatabasePropertiesUpdated(event *sharedEvents.DatabasePropertiesUpdated) error {
	dto := sharedDTO.DatabasePropertiesUpdatedDTOFromDomain(event)
	return wt.publishDTO(sharedEvents.DatabasePropertiesUpdatedTopic, dto)
}

// publishDTO is a helper method to publish DTOs to Watermill
func (wt *WebhookTriage) publishDTO(topic string, dto interface{}) error {
	eventBytes, err := json.Marshal(dto)
	if err != nil {
		wt.logger.Printf("Failed to marshal event DTO: %v", err)
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), eventBytes)

	if err := wt.publisher.Publish(topic, msg); err != nil {
		wt.logger.Printf("Failed to publish event to topic %s: %v", topic, err)
		return err
	}

	wt.logger.Printf("Successfully published event to topic: %s", topic)
	return nil
}

// Handler returns a Watermill consumer handler function
func (wt *WebhookTriage) Handler() message.NoPublishHandlerFunc {
	return func(msg *message.Message) error {
		return wt.ProcessWebhook(msg)
	}
}
