package events

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	projects "src/internal/modules/projects/domain"
	sharedEvents "src/internal/modules/shared/domain/events"
	sharedDTO "src/internal/modules/shared/interfaces/dto"
	tasks "src/internal/modules/tasks/domain"
)

type stubPublisher struct {
	topics   []string
	messages []*message.Message
}

func (p *stubPublisher) Publish(topic string, messages ...*message.Message) error {
	p.topics = append(p.topics, topic)
	p.messages = append(p.messages, messages...)
	return nil
}

func (p *stubPublisher) Close() error {
	return nil
}

type fakeTasksRepo struct {
	findByNotionIDs func(ctx context.Context, notionDatabaseID, notionPageID string) (*tasks.Task, error)
	upsert          func(ctx context.Context, task *tasks.Task) error
}

func (r *fakeTasksRepo) Save(ctx context.Context, task *tasks.Task) error {
	return nil
}

func (r *fakeTasksRepo) Update(ctx context.Context, task *tasks.Task) error {
	return nil
}

func (r *fakeTasksRepo) Upsert(ctx context.Context, task *tasks.Task) error {
	if r.upsert != nil {
		return r.upsert(ctx, task)
	}
	return nil
}

func (r *fakeTasksRepo) FindByID(ctx context.Context, id uuid.UUID) (*tasks.Task, error) {
	return nil, tasks.ErrTaskNotFound
}

func (r *fakeTasksRepo) FindByPublicID(ctx context.Context, publicID string) (*tasks.Task, error) {
	return nil, tasks.ErrTaskNotFound
}

func (r *fakeTasksRepo) FindByNotionIDs(ctx context.Context, notionDatabaseID, notionPageID string) (*tasks.Task, error) {
	if r.findByNotionIDs != nil {
		return r.findByNotionIDs(ctx, notionDatabaseID, notionPageID)
	}
	return nil, tasks.ErrTaskNotFound
}

func (r *fakeTasksRepo) WithTransaction(ctx context.Context, fn func(ctx context.Context, repo tasks.Repository) error) error {
	return fn(ctx, r)
}

type fakeProjectsRepo struct {
	findByNotionDatabaseID func(ctx context.Context, notionDatabaseID string) (*projects.Project, error)
}

func (r *fakeProjectsRepo) Save(ctx context.Context, project *projects.Project) error {
	return nil
}

func (r *fakeProjectsRepo) FindByID(ctx context.Context, id uuid.UUID) (*projects.Project, error) {
	return nil, projects.ErrProjectNotFound
}

func (r *fakeProjectsRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*projects.Project, error) {
	return nil, nil
}

func (r *fakeProjectsRepo) FindByNotionDatabaseID(ctx context.Context, notionDatabaseID string) (*projects.Project, error) {
	if r.findByNotionDatabaseID != nil {
		return r.findByNotionDatabaseID(ctx, notionDatabaseID)
	}
	return nil, projects.ErrProjectNotFound
}

func (r *fakeProjectsRepo) Update(ctx context.Context, project *projects.Project) error {
	return nil
}

func (r *fakeProjectsRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

var _ = Describe("Webhook triage", func() {
	Describe("ID extraction", func() {
		It("uses entity id for page events", func() {
			payload := NotionWebhookPayload{
				Entity: NotionWebhookEntity{ID: "page_123", Type: "page"},
			}
			pageID, ok := getPageID(payload)
			Expect(ok).To(BeTrue())
			Expect(pageID).To(Equal("page_123"))
		})

		It("falls back to data.page_id or data.id for page id", func() {
			payload := NotionWebhookPayload{
				Data: map[string]interface{}{"page_id": "page_456"},
			}
			pageID, ok := getPageID(payload)
			Expect(ok).To(BeTrue())
			Expect(pageID).To(Equal("page_456"))

			payload = NotionWebhookPayload{
				Data: map[string]interface{}{"id": "page_789"},
			}
			pageID, ok = getPageID(payload)
			Expect(ok).To(BeTrue())
			Expect(pageID).To(Equal("page_789"))
		})

		It("uses entity id for database and data_source events", func() {
			payload := NotionWebhookPayload{
				Entity: NotionWebhookEntity{ID: "db_123", Type: "database"},
			}
			databaseID, ok := getDatabaseID(payload)
			Expect(ok).To(BeTrue())
			Expect(databaseID).To(Equal("db_123"))

			payload = NotionWebhookPayload{
				Entity: NotionWebhookEntity{ID: "ds_123", Type: "data_source"},
			}
			databaseID, ok = getDatabaseID(payload)
			Expect(ok).To(BeTrue())
			Expect(databaseID).To(Equal("ds_123"))
		})

		It("falls back to data.parent for database id", func() {
			payload := NotionWebhookPayload{
				Data: map[string]interface{}{
					"parent": map[string]interface{}{
						"id":   "db_parent",
						"type": "database",
					},
				},
			}
			databaseID, ok := getDatabaseID(payload)
			Expect(ok).To(BeTrue())
			Expect(databaseID).To(Equal("db_parent"))

			payload = NotionWebhookPayload{
				Data: map[string]interface{}{
					"parent": map[string]interface{}{
						"database_id": "db_fallback",
					},
				},
			}
			databaseID, ok = getDatabaseID(payload)
			Expect(ok).To(BeTrue())
			Expect(databaseID).To(Equal("db_fallback"))
		})
	})

	Describe("routing and publishing", func() {
		It("publishes a TaskPropertiesUpdated event for page.content_updated", func() {
			projectID := uuid.New()
			var upsertedTask *tasks.Task

			tasksRepo := &fakeTasksRepo{
				findByNotionIDs: func(ctx context.Context, notionDatabaseID, notionPageID string) (*tasks.Task, error) {
					return nil, tasks.ErrTaskNotFound
				},
				upsert: func(ctx context.Context, task *tasks.Task) error {
					upsertedTask = task
					return nil
				},
			}
			projectsRepo := &fakeProjectsRepo{
				findByNotionDatabaseID: func(ctx context.Context, notionDatabaseID string) (*projects.Project, error) {
					return &projects.Project{
						ID:               projectID,
						NotionDatabaseID: notionDatabaseID,
					}, nil
				},
			}
			publisher := &stubPublisher{}
			logger := log.New(io.Discard, "", 0)

			triage := NewWebhookTriage(publisher, logger, tasksRepo, projectsRepo, nil, nil)
			payload := NotionWebhookPayload{
				Type: "page.content_updated",
				Entity: NotionWebhookEntity{
					ID:   "page_123",
					Type: "page",
				},
				Data: map[string]interface{}{
					"parent": map[string]interface{}{
						"id":   "db_123",
						"type": "database",
					},
				},
			}

			rawPayload, err := json.Marshal(payload)
			Expect(err).ToNot(HaveOccurred())

			envelope, err := json.Marshal(sharedEvents.NotionWebhookReceived{Payload: rawPayload})
			Expect(err).ToNot(HaveOccurred())

			msg := message.NewMessage("msg_1", envelope)
			Expect(triage.ProcessWebhook(msg)).To(Succeed())

			Expect(upsertedTask).ToNot(BeNil())
			Expect(publisher.topics).To(ConsistOf(sharedEvents.TaskPropertiesUpdatedTopic))
			Expect(publisher.messages).To(HaveLen(1))

			var dto sharedDTO.TaskPropertiesUpdatedDTO
			Expect(json.Unmarshal(publisher.messages[0].Payload, &dto)).To(Succeed())
			Expect(dto.TaskID).To(Equal(upsertedTask.ID))
			Expect(dto.DatabaseID).To(Equal("db_123"))
			Expect(dto.PageID).To(Equal("page_123"))
		})

		It("publishes a TaskCreated event for page.created", func() {
			projectID := uuid.New()
			var upsertedTask *tasks.Task

			tasksRepo := &fakeTasksRepo{
				findByNotionIDs: func(ctx context.Context, notionDatabaseID, notionPageID string) (*tasks.Task, error) {
					return nil, tasks.ErrTaskNotFound
				},
				upsert: func(ctx context.Context, task *tasks.Task) error {
					upsertedTask = task
					return nil
				},
			}
			projectsRepo := &fakeProjectsRepo{
				findByNotionDatabaseID: func(ctx context.Context, notionDatabaseID string) (*projects.Project, error) {
					return &projects.Project{
						ID:               projectID,
						NotionDatabaseID: notionDatabaseID,
					}, nil
				},
			}
			publisher := &stubPublisher{}
			logger := log.New(io.Discard, "", 0)

			triage := NewWebhookTriage(publisher, logger, tasksRepo, projectsRepo, nil, nil)
			payload := NotionWebhookPayload{
				Type: "page.created",
				Entity: NotionWebhookEntity{
					ID:   "page_created",
					Type: "page",
				},
				Data: map[string]interface{}{
					"parent": map[string]interface{}{
						"id":   "db_created",
						"type": "database",
					},
				},
			}

			rawPayload, err := json.Marshal(payload)
			Expect(err).ToNot(HaveOccurred())

			envelope, err := json.Marshal(sharedEvents.NotionWebhookReceived{Payload: rawPayload})
			Expect(err).ToNot(HaveOccurred())

			msg := message.NewMessage("msg_3", envelope)
			Expect(triage.ProcessWebhook(msg)).To(Succeed())

			Expect(upsertedTask).ToNot(BeNil())
			Expect(publisher.topics).To(ConsistOf(sharedEvents.TaskCreatedTopic))
			Expect(publisher.messages).To(HaveLen(1))

			var dto sharedDTO.TaskCreatedDTO
			Expect(json.Unmarshal(publisher.messages[0].Payload, &dto)).To(Succeed())
			Expect(dto.TaskID).To(Equal(upsertedTask.ID))
			Expect(dto.DatabaseID).To(Equal("db_created"))
			Expect(dto.PageID).To(Equal("page_created"))
		})

		It("publishes a TaskDeleted event for page.deleted when a task exists", func() {
			existingTask := &tasks.Task{
				ID:               uuid.New(),
				NotionDatabaseID: "db_deleted",
				NotionPageID:     "page_deleted",
			}

			tasksRepo := &fakeTasksRepo{
				findByNotionIDs: func(ctx context.Context, notionDatabaseID, notionPageID string) (*tasks.Task, error) {
					return existingTask, nil
				},
			}
			projectsRepo := &fakeProjectsRepo{}
			publisher := &stubPublisher{}
			logger := log.New(io.Discard, "", 0)

			triage := NewWebhookTriage(publisher, logger, tasksRepo, projectsRepo, nil, nil)
			payload := NotionWebhookPayload{
				Type: "page.deleted",
				Entity: NotionWebhookEntity{
					ID:   "page_deleted",
					Type: "page",
				},
				Data: map[string]interface{}{
					"parent": map[string]interface{}{
						"id":   "db_deleted",
						"type": "database",
					},
				},
			}

			rawPayload, err := json.Marshal(payload)
			Expect(err).ToNot(HaveOccurred())

			envelope, err := json.Marshal(sharedEvents.NotionWebhookReceived{Payload: rawPayload})
			Expect(err).ToNot(HaveOccurred())

			msg := message.NewMessage("msg_4", envelope)
			Expect(triage.ProcessWebhook(msg)).To(Succeed())

			Expect(publisher.topics).To(ConsistOf(sharedEvents.TaskDeletedTopic))
			Expect(publisher.messages).To(HaveLen(1))

			var dto sharedDTO.TaskDeletedDTO
			Expect(json.Unmarshal(publisher.messages[0].Payload, &dto)).To(Succeed())
			Expect(dto.TaskID).To(Equal(existingTask.ID))
			Expect(dto.DatabaseID).To(Equal("db_deleted"))
			Expect(dto.PageID).To(Equal("page_deleted"))
		})

		It("publishes a DatabasePropertiesUpdated event for database.schema_updated", func() {
			publisher := &stubPublisher{}
			logger := log.New(io.Discard, "", 0)

			triage := NewWebhookTriage(publisher, logger, nil, nil, nil, nil)
			payload := NotionWebhookPayload{
				Type: "database.schema_updated",
				Entity: NotionWebhookEntity{
					ID:   "db_schema",
					Type: "database",
				},
				Data: map[string]interface{}{
					"updated_properties": []interface{}{"prop_1"},
				},
			}

			rawPayload, err := json.Marshal(payload)
			Expect(err).ToNot(HaveOccurred())

			envelope, err := json.Marshal(sharedEvents.NotionWebhookReceived{Payload: rawPayload})
			Expect(err).ToNot(HaveOccurred())

			msg := message.NewMessage("msg_5", envelope)
			Expect(triage.ProcessWebhook(msg)).To(Succeed())
			Expect(publisher.topics).To(ConsistOf(sharedEvents.DatabasePropertiesUpdatedTopic))
			Expect(publisher.messages).To(HaveLen(1))

			var dto sharedDTO.DatabasePropertiesUpdatedDTO
			Expect(json.Unmarshal(publisher.messages[0].Payload, &dto)).To(Succeed())
			Expect(dto.DatabaseID).To(Equal("db_schema"))
			Expect(dto.Changes).To(HaveKey("updated_properties"))
		})

		It("ignores known but unhandled event types", func() {
			publisher := &stubPublisher{}
			logger := log.New(io.Discard, "", 0)

			triage := NewWebhookTriage(publisher, logger, nil, nil, nil, nil)
			payload := NotionWebhookPayload{
				Type: "page.locked",
				Entity: NotionWebhookEntity{
					ID:   "page_123",
					Type: "page",
				},
			}

			rawPayload, err := json.Marshal(payload)
			Expect(err).ToNot(HaveOccurred())

			envelope, err := json.Marshal(sharedEvents.NotionWebhookReceived{Payload: rawPayload})
			Expect(err).ToNot(HaveOccurred())

			msg := message.NewMessage("msg_2", envelope)
			Expect(triage.ProcessWebhook(msg)).To(Succeed())
			Expect(publisher.messages).To(BeEmpty())
			Expect(publisher.topics).To(BeEmpty())
		})
	})
})
