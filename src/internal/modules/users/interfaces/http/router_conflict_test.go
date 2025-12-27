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

	shared "src/internal/modules/shared/domain"
	"src/internal/modules/users/domain"
	userhttp "src/internal/modules/users/interfaces/http"
)

type conflictUserRepo struct {
	existing bool
}

func (c conflictUserRepo) Create(_ context.Context, user domain.User) (domain.User, error) {
	if c.existing {
		return domain.User{}, domain.ErrUserAlreadyExists
	}
	return user, nil
}

func (c conflictUserRepo) GetByID(_ context.Context, id string) (domain.User, error) {
	return domain.User{}, domain.ErrUserNotFound
}

func (c conflictUserRepo) GetByEmail(_ context.Context, email string) (domain.User, error) {
	if c.existing {
		return domain.User{Email: email, ID: uuid.New()}, nil
	}
	return domain.User{}, domain.ErrUserNotFound
}

func (c conflictUserRepo) Update(_ context.Context, user domain.User) (domain.User, error) {
	return user, nil
}

func (c conflictUserRepo) Delete(_ context.Context, id string) error { return nil }
func (c conflictUserRepo) List(_ context.Context, offset, limit int) ([]domain.User, error) {
	return nil, nil
}

type fixedClock struct{ now time.Time }

func (f fixedClock) Now() time.Time { return f.now }

var _ = Describe("Users Router", func() {
	It("returns 409 when user already exists", func() {
		repo := conflictUserRepo{existing: true}
		router := userhttp.NewRouter(repo, shared.NewNoopTransactionManager(), fixedIDGen{}, fixedClock{now: time.Now()})

		body := bytes.NewBufferString(`{"email":"dup@example.com","name":"Dup User"}`)
		req := httptest.NewRequest(http.MethodPost, "/", body)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		Expect(rec.Code).To(Equal(http.StatusConflict))
		Expect(rec.Body.String()).To(ContainSubstring("user already exists"))
	})
})
