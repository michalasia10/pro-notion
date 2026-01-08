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
	pageContentUpdatedEvent    = "page.content_updated"
	pageCreatedEvent           = "page.created"
	pageDeletedEvent           = "page.deleted"
	pageLockedEvent            = "page.locked"
	pageMovedEvent             = "page.moved"
	pagePropertiesUpdatedEvent = "page.properties_updated"
	pageUndeletedEvent         = "page.undeleted"
	pageUnlockedEvent          = "page.unlocked"

	databaseContentUpdatedEvent = "database.content_updated"
	databaseCreatedEvent        = "database.created"
	databaseDeletedEvent        = "database.deleted"
	databaseMovedEvent          = "database.moved"
	databaseSchemaUpdatedEvent  = "database.schema_updated"
	databaseUndeletedEvent      = "database.undeleted"

	dataSourceContentUpdatedEvent = "data_source.content_updated"
	dataSourceCreatedEvent        = "data_source.created"
	dataSourceDeletedEvent        = "data_source.deleted"
	dataSourceMovedEvent          = "data_source.moved"
	dataSourceSchemaUpdatedEvent  = "data_source.schema_updated"
	dataSourceUndeletedEvent      = "data_source.undeleted"

	commentCreatedEvent = "comment.created"
	commentDeletedEvent = "comment.deleted"
	commentUpdatedEvent = "comment.updated"
)

// NotionWebhookPayload represents the structure of a webhook payload from Notion
type NotionWebhookPayload struct {
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp string                 `json:"timestamp"`
	Entity    NotionWebhookEntity    `json:"entity"`
}

type NotionWebhookEntity struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

var knownWebhookEventTypes = map[string]struct{}{
	pageContentUpdatedEvent:    {},
	pageCreatedEvent:           {},
	pageDeletedEvent:           {},
	pageLockedEvent:            {},
	pageMovedEvent:             {},
	pagePropertiesUpdatedEvent: {},
	pageUndeletedEvent:         {},
	pageUnlockedEvent:          {},
	databaseContentUpdatedEvent: {},
	databaseCreatedEvent:        {},
	databaseDeletedEvent:        {},
	databaseMovedEvent:          {},
	databaseSchemaUpdatedEvent:  {},
	databaseUndeletedEvent:      {},
	dataSourceContentUpdatedEvent: {},
	dataSourceCreatedEvent:        {},
	dataSourceDeletedEvent:        {},
	dataSourceMovedEvent:          {},
	dataSourceSchemaUpdatedEvent:  {},
	dataSourceUndeletedEvent:      {},
	commentCreatedEvent: {},
	commentDeletedEvent: {},
	commentUpdatedEvent: {},
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
	case pageContentUpdatedEvent:
		return wt.handlePageUpdated(ctx, notionPayload)
	case pagePropertiesUpdatedEvent:
		return wt.handlePageUpdated(ctx, notionPayload)
	case pageCreatedEvent:
		return wt.handlePageCreated(ctx, notionPayload)
	case pageDeletedEvent:
		return wt.handlePageDeleted(ctx, notionPayload)
	case databaseSchemaUpdatedEvent:
		return wt.handleDatabaseUpdated(notionPayload)
	case dataSourceSchemaUpdatedEvent:
		return wt.handleDatabaseUpdated(notionPayload)
	default:
		if _, ok := knownWebhookEventTypes[notionPayload.Type]; ok {
			wt.logger.Printf("Ignoring webhook type: %s", notionPayload.Type)
			return nil
		}
		wt.logger.Printf("Unknown webhook type: %s", notionPayload.Type)
		return nil
	}
}

func (wt *WebhookTriage) logEvent(eventType, message string) {
	wt.logger.Printf("%s| %s", eventType, message)
}

func (wt *WebhookTriage) handlePageUpdated(ctx context.Context, payload NotionWebhookPayload) error {
	eventType := payload.Type
	if eventType == "" {
		eventType = pagePropertiesUpdatedEvent
	}
	pageID, ok := getPageID(payload)
	if !ok {
		wt.logEvent(eventType, "Missing page ID")
		return nil
	}

	databaseID, ok := getDatabaseID(payload)
	if !ok {
		if isNonDatabaseParent(payload) {
			wt.logEvent(eventType, "Ignoring page event with non-database parent")
			return nil
		}
		wt.logEvent(eventType, "Missing database ID")
		return nil
	}

	taskID, created, err := wt.resolveTaskID(ctx, databaseID, pageID, true)
	if err != nil {
		wt.logEvent(eventType, fmt.Sprintf("Failed to resolve task for (db=%s page=%s): %v", databaseID, pageID, err))
		return err
	}

	if created {
		wt.logEvent(eventType, fmt.Sprintf("Created new task mapping db=%s page=%s task=%s", databaseID, pageID, taskID))
	}

	event := sharedEvents.NewTaskPropertiesUpdated(taskID, databaseID, pageID, payload.Data)
	return wt.publishTaskPropertiesUpdated(event)
}

func (wt *WebhookTriage) handlePageCreated(ctx context.Context, payload NotionWebhookPayload) error {
	pageID, ok := getPageID(payload)
	if !ok {
		wt.logEvent(pageCreatedEvent, "Missing page ID")
		return nil
	}

	databaseID, ok := getDatabaseID(payload)
	if !ok {
		if isNonDatabaseParent(payload) {
			wt.logEvent(pageCreatedEvent, "Ignoring page event with non-database parent")
			return nil
		}
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
	pageID, ok := getPageID(payload)
	if !ok {
		wt.logEvent(pageDeletedEvent, "Missing page ID")
		return nil
	}

	databaseID, ok := getDatabaseID(payload)
	if !ok {
		if isNonDatabaseParent(payload) {
			wt.logEvent(pageDeletedEvent, "Ignoring page event with non-database parent")
			return nil
		}
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
	eventType := payload.Type
	if eventType == "" {
		eventType = databaseSchemaUpdatedEvent
	}
	databaseID, ok := getDatabaseID(payload)
	if !ok {
		wt.logEvent(eventType, "Missing database ID")
		return nil
	}

	pageID := ""
	if payload.Data != nil {
		if id, ok := payload.Data["page_id"].(string); ok {
			pageID = id
		}
	}

	event := sharedEvents.NewDatabasePropertiesUpdated(databaseID, pageID, payload.Data)
	return wt.publishDatabasePropertiesUpdated(event)
}

func getPageID(payload NotionWebhookPayload) (string, bool) {
	if payload.Entity.Type == "page" && payload.Entity.ID != "" {
		return payload.Entity.ID, true
	}
	if payload.Data == nil {
		return "", false
	}
	if id, ok := payload.Data["page_id"].(string); ok && id != "" {
		return id, true
	}
	if id, ok := payload.Data["id"].(string); ok && id != "" {
		return id, true
	}
	return "", false
}

func getDatabaseID(payload NotionWebhookPayload) (string, bool) {
	if (payload.Entity.Type == "database" || payload.Entity.Type == "data_source") && payload.Entity.ID != "" {
		return payload.Entity.ID, true
	}
	if payload.Data == nil {
		return "", false
	}
	parent, ok := payload.Data["parent"].(map[string]interface{})
	if !ok {
		return "", false
	}
	if parentType, ok := parent["type"].(string); ok && parentType == "database" {
		if id, ok := parent["id"].(string); ok && id != "" {
			return id, true
		}
	}
	if id, ok := parent["database_id"].(string); ok && id != "" {
		return id, true
	}
	return "", false
}

func isNonDatabaseParent(payload NotionWebhookPayload) bool {
	if payload.Entity.Type != "page" {
		return false
	}
	if payload.Data == nil {
		return false
	}
	parent, ok := payload.Data["parent"].(map[string]interface{})
	if !ok {
		return false
	}
	parentType, ok := parent["type"].(string)
	if !ok {
		return false
	}
	return parentType != "database"
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
