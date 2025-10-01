package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"

	projects "src/internal/modules/projects/domain"
	shared "src/internal/modules/shared/domain"
	sharedEvents "src/internal/modules/shared/domain/events"
	sharedDTO "src/internal/modules/shared/interfaces/dto"
	tasks "src/internal/modules/tasks/domain"
)

var (
	pageUpdatedEvent     = "page.updated"
	pageCreatedEvent     = "page.created"
	pageDeletedEvent     = "page.deleted"
	databaseUpdatedEvent = "database.updated"
)

// NotionWebhookPayload represents the structure of a webhook payload from Notion
type NotionWebhookPayload struct {
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// WebhookTriage processes raw Notion webhook events and publishes specific domain events
type WebhookTriage struct {
	publisher    message.Publisher
	logger       *log.Logger
	tasksRepo    tasks.Repository
	projectsRepo projects.Repository
	idGenerator  shared.IDGenerator
	clock        func() time.Time
}

// NewWebhookTriage creates a new webhook triage processor
func NewWebhookTriage(
	publisher message.Publisher,
	logger *log.Logger,
	tasksRepo tasks.Repository,
	projectsRepo projects.Repository,
	idGenerator shared.IDGenerator,
	clock func() time.Time,
) *WebhookTriage {
	if idGenerator == nil {
		idGenerator = shared.NewUUIDGenerator()
	}
	if clock == nil {
		clock = time.Now
	}

	return &WebhookTriage{
		publisher:    publisher,
		logger:       logger,
		tasksRepo:    tasksRepo,
		projectsRepo: projectsRepo,
		idGenerator:  idGenerator,
		clock:        clock,
	}
}

// ProcessWebhook handles incoming webhook messages from Watermill
func (wt *WebhookTriage) ProcessWebhook(msg *message.Message) error {
	defer msg.Ack()

	var webhookEvent sharedEvents.NotionWebhookReceived
	if err := json.Unmarshal(msg.Payload, &webhookEvent); err != nil {
		wt.logger.Printf("Failed to unmarshal webhook event: %v", err)
		return err
	}

	var notionPayload NotionWebhookPayload
	if err := json.Unmarshal(webhookEvent.Payload, &notionPayload); err != nil {
		wt.logger.Printf("Failed to unmarshal Notion payload: %v", err)
		return err
	}

	ctx := msg.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	switch notionPayload.Type {
	case pageUpdatedEvent:
		return wt.handlePageUpdated(ctx, notionPayload)
	case pageCreatedEvent:
		return wt.handlePageCreated(ctx, notionPayload)
	case pageDeletedEvent:
		return wt.handlePageDeleted(ctx, notionPayload)
	case databaseUpdatedEvent:
		return wt.handleDatabaseUpdated(notionPayload)
	default:
		wt.logger.Printf("Unknown webhook type: %s", notionPayload.Type)
		return nil
	}
}

func (wt *WebhookTriage) logEvent(eventType, message string) {
	wt.logger.Printf("%s| %s", eventType, message)
}

func (wt *WebhookTriage) handlePageUpdated(ctx context.Context, payload NotionWebhookPayload) error {
	pageID, ok := payload.Data["id"].(string)
	if !ok {
		wt.logEvent(pageUpdatedEvent, "Missing page ID")
		return nil
	}

	databaseID, ok := payload.Data["database_id"].(string)
	if !ok {
		wt.logEvent(pageUpdatedEvent, "Missing database ID")
		return nil
	}

	taskID, created, err := wt.resolveTaskID(ctx, databaseID, pageID, true)
	if err != nil {
		wt.logEvent(pageUpdatedEvent, fmt.Sprintf("Failed to resolve task for (db=%s page=%s): %v", databaseID, pageID, err))
		return err
	}

	if created {
		wt.logEvent(pageUpdatedEvent, fmt.Sprintf("Created new task mapping db=%s page=%s task=%s", databaseID, pageID, taskID))
	}

	event := sharedEvents.NewTaskPropertiesUpdated(taskID, databaseID, pageID, payload.Data)
	return wt.publishTaskPropertiesUpdated(event)
}

func (wt *WebhookTriage) handlePageCreated(ctx context.Context, payload NotionWebhookPayload) error {
	pageID, ok := payload.Data["id"].(string)
	if !ok {
		wt.logEvent(pageCreatedEvent, "Missing page ID")
		return nil
	}

	databaseID, ok := payload.Data["database_id"].(string)
	if !ok {
		wt.logEvent(pageCreatedEvent, "Missing database ID")
		return nil
	}

	taskID, created, err := wt.resolveTaskID(ctx, databaseID, pageID, true)
	if err != nil {
		wt.logEvent(pageCreatedEvent, fmt.Sprintf("Failed to resolve task for (db=%s page=%s): %v", databaseID, pageID, err))
		return err
	}

	if created {
		wt.logEvent(pageCreatedEvent, fmt.Sprintf("Created new task mapping db=%s page=%s task=%s", databaseID, pageID, taskID))
	}

	event := sharedEvents.NewTaskCreated(taskID, databaseID, pageID, payload.Data)
	return wt.publishTaskCreated(event)
}

func (wt *WebhookTriage) handlePageDeleted(ctx context.Context, payload NotionWebhookPayload) error {
	pageID, ok := payload.Data["id"].(string)
	if !ok {
		wt.logEvent(pageDeletedEvent, "Missing page ID")
		return nil
	}

	databaseID, ok := payload.Data["database_id"].(string)
	if !ok {
		wt.logEvent(pageDeletedEvent, "Missing database ID")
		return nil
	}

	taskID, _, err := wt.resolveTaskID(ctx, databaseID, pageID, false)
	if err != nil {
		if errors.Is(err, tasks.ErrTaskNotFound) {
			wt.logEvent(pageDeletedEvent, fmt.Sprintf("No matching task found for deletion (db=%s page=%s)", databaseID, pageID))
			return nil
		}
		wt.logEvent(pageDeletedEvent, fmt.Sprintf("Failed to resolve task for (db=%s page=%s): %v", databaseID, pageID, err))
		return err
	}

	event := sharedEvents.NewTaskDeleted(taskID, databaseID, pageID)
	return wt.publishTaskDeleted(event)
}

func (wt *WebhookTriage) handleDatabaseUpdated(payload NotionWebhookPayload) error {
	databaseID, ok := payload.Data["id"].(string)
	if !ok {
		wt.logEvent(databaseUpdatedEvent, "Missing database ID")
		return nil
	}

	pageID, ok := payload.Data["page_id"].(string)
	if !ok {
		pageID = ""
	}

	event := sharedEvents.NewDatabasePropertiesUpdated(databaseID, pageID, payload.Data)
	return wt.publishDatabasePropertiesUpdated(event)
}

func (wt *WebhookTriage) resolveTaskID(ctx context.Context, databaseID, pageID string, createIfMissing bool) (uuid.UUID, bool, error) {
	task, err := wt.tasksRepo.FindByNotionIDs(ctx, databaseID, pageID)
	if err == nil {
		return task.ID, false, nil
	}

	if !errors.Is(err, tasks.ErrTaskNotFound) {
		return uuid.Nil, false, err
	}

	if !createIfMissing {
		return uuid.Nil, false, tasks.ErrTaskNotFound
	}

	project, err := wt.projectsRepo.FindByNotionDatabaseID(ctx, databaseID)
	if err != nil {
		return uuid.Nil, false, err
	}

	now := wt.clock()
	newTask, err := tasks.NewTask(project.ID, databaseID, pageID, now, wt.idGenerator)
	if err != nil {
		return uuid.Nil, false, err
	}

	if err := wt.tasksRepo.Upsert(ctx, newTask); err != nil {
		return uuid.Nil, false, err
	}

	return newTask.ID, true, nil
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
		wt.logEvent(topic, fmt.Sprintf("Failed to marshal event DTO: %v", err))
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), eventBytes)

	if err := wt.publisher.Publish(topic, msg); err != nil {
		wt.logEvent(topic, fmt.Sprintf("Failed to publish event: %v", err))
		return err
	}

	wt.logEvent(topic, "Successfully published event")
	return nil
}

// Handler returns a Watermill consumer handler function
func (wt *WebhookTriage) Handler() message.NoPublishHandlerFunc {
	return func(msg *message.Message) error {
		return wt.ProcessWebhook(msg)
	}
}
