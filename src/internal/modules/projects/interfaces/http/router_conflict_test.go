package http_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"src/internal/modules/projects/domain"
	projecthttp "src/internal/modules/projects/interfaces/http"
	shared "src/internal/modules/shared/domain"
	"src/internal/pkg/middleware"
)

type conflictProjectRepo struct {
	existing bool
}

func (c *conflictProjectRepo) Save(ctx context.Context, project *domain.Project) error {
	if c.existing {
		return domain.ErrProjectAlreadyExists
	}
	return nil
}

func (c *conflictProjectRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	return nil, domain.ErrProjectNotFound
}

func (c *conflictProjectRepo) FindByPublicID(ctx context.Context, publicID string) (*domain.Project, error) {
	return nil, domain.ErrProjectNotFound
}

func (c *conflictProjectRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Project, error) {
	return nil, nil
}

func (c *conflictProjectRepo) FindByNotionDatabaseID(ctx context.Context, notionDatabaseID string) (*domain.Project, error) {
	if c.existing {
		return &domain.Project{ID: uuid.New(), NotionDatabaseID: notionDatabaseID}, nil
	}
	return nil, domain.ErrProjectNotFound
}

func (c *conflictProjectRepo) Update(ctx context.Context, project *domain.Project) error { return nil }
func (c *conflictProjectRepo) Delete(ctx context.Context, id uuid.UUID) error            { return nil }

type fixedClock struct{ now time.Time }

func (f fixedClock) Now() time.Time { return f.now }

type fixedIDGen struct{}

func (fixedIDGen) NewID(prefix string) string { return prefix + "_id" }

var _ = Describe("Projects Router", func() {
	It("returns 409 when project already exists for database", func() {
		repo := &conflictProjectRepo{existing: true}
		router := projecthttp.NewRouter(repo, shared.NewNoopTransactionManager(), fixedIDGen{}, fixedClock{now: time.Now()})

		body := bytes.NewBufferString(`{"notion_database_id":"db_1","notion_webhook_secret":"secret"}`)
		req := httptest.NewRequest(http.MethodPost, "/", body)
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New()))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		Expect(rec.Code).To(Equal(http.StatusConflict))
		Expect(rec.Body.String()).To(ContainSubstring("project already exists"))
	})
})
