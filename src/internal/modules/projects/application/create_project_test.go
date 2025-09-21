package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"src/internal/modules/projects/application"
	"src/internal/modules/projects/domain"
	shared "src/internal/modules/shared/domain"
)

// Mock implementations for testing
type mockProjectRepository struct {
	projects map[uuid.UUID]*domain.Project
}

func newMockProjectRepository() *mockProjectRepository {
	return &mockProjectRepository{
		projects: make(map[uuid.UUID]*domain.Project),
	}
}

func (m *mockProjectRepository) Save(ctx context.Context, project *domain.Project) error {
	m.projects[project.ID] = project
	return nil
}

func (m *mockProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	project, exists := m.projects[id]
	if !exists {
		return nil, domain.ErrProjectNotFound
	}
	return project, nil
}

func (m *mockProjectRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Project, error) {
	var projects []*domain.Project
	for _, p := range m.projects {
		if p.UserID == userID {
			projects = append(projects, p)
		}
	}
	return projects, nil
}

func (m *mockProjectRepository) FindByNotionDatabaseID(ctx context.Context, notionDatabaseID string) (*domain.Project, error) {
	for _, p := range m.projects {
		if p.NotionDatabaseID == notionDatabaseID {
			return p, nil
		}
	}
	return nil, domain.ErrProjectNotFound
}

func (m *mockProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	m.projects[project.ID] = project
	return nil
}

func (m *mockProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.projects, id)
	return nil
}

type mockClock struct {
	now time.Time
}

func (m *mockClock) Now() time.Time {
	return m.now
}

type mockIDGenerator struct {
	counter int
}

func (m *mockIDGenerator) NewID(prefix string) string {
	m.counter++
	return prefix + "_test_" + string(rune(m.counter+'0'))
}

type mockTransactionManager struct {
	shouldFail bool
}

func (m *mockTransactionManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if m.shouldFail {
		return errors.New("transaction failed")
	}
	return fn(ctx)
}

var _ = Describe("CreateProjectUseCase", func() {
	var (
		repo  domain.Repository
		idGen shared.IDGenerator
		clock shared.Clock
		txMgr shared.TransactionManager
		uc    *application.CreateProjectUseCase
		ctx   context.Context
	)

	BeforeEach(func() {
		repo = newMockProjectRepository()
		idGen = &mockIDGenerator{}
		clock = &mockClock{now: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)}
		txMgr = &mockTransactionManager{}
		uc = application.NewCreateProjectUseCase(repo, idGen, clock, txMgr)
		ctx = context.Background()
	})

	Describe("Execute", func() {
		It("should create a project successfully", func() {
			req := application.CreateProjectRequest{
				UserID:              uuid.New(),
				NotionDatabaseID:    "database_123",
				NotionWebhookSecret: "secret_123",
			}

			resp, err := uc.Execute(ctx, req)

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Project.ID).ToNot(BeNil())
			Expect(resp.Project.UserID).To(Equal(req.UserID))
			Expect(resp.Project.NotionDatabaseID).To(Equal(req.NotionDatabaseID))
			Expect(resp.Project.NotionWebhookSecret).To(Equal(req.NotionWebhookSecret))
		})

		It("should return error when project already exists for the database", func() {
			// Create first project
			req1 := application.CreateProjectRequest{
				UserID:              uuid.New(),
				NotionDatabaseID:    "database_123",
				NotionWebhookSecret: "secret_123",
			}
			_, err := uc.Execute(ctx, req1)
			Expect(err).ToNot(HaveOccurred())

			// Try to create second project with same database ID
			req2 := application.CreateProjectRequest{
				UserID:              uuid.New(),
				NotionDatabaseID:    "database_123",
				NotionWebhookSecret: "secret_456",
			}
			_, err = uc.Execute(ctx, req2)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(domain.ErrProjectNotFound)) // This should probably be a different error
		})

		It("should return error when transaction fails", func() {
			txMgr := &mockTransactionManager{shouldFail: true}
			uc := application.NewCreateProjectUseCase(repo, idGen, clock, txMgr)

			req := application.CreateProjectRequest{
				UserID:              uuid.New(),
				NotionDatabaseID:    "database_123",
				NotionWebhookSecret: "secret_123",
			}

			_, err := uc.Execute(ctx, req)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("transaction failed"))
		})
	})
})

func TestCreateProject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateProject Application Suite")
}
