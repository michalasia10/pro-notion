package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"src/internal/config"
	shared "src/internal/modules/shared/domain"
	"src/internal/modules/users/domain"
	userhttp "src/internal/modules/users/interfaces/http"
	"src/internal/pkg/notion"
)

type dummyUserRepo struct{}

func (d dummyUserRepo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	return user, nil
}
func (d dummyUserRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	return domain.User{}, domain.ErrUserNotFound
}
func (d dummyUserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return domain.User{}, domain.ErrUserNotFound
}
func (d dummyUserRepo) Update(ctx context.Context, user domain.User) (domain.User, error) {
	return user, nil
}
func (d dummyUserRepo) Delete(ctx context.Context, id string) error { return nil }
func (d dummyUserRepo) List(ctx context.Context, offset, limit int) ([]domain.User, error) {
	return []domain.User{}, nil
}

type noopTx struct{}

func (n noopTx) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type fixedIDGen struct{}

func (f fixedIDGen) NewID(prefix string) string { return prefix + "_id" }

var _ = Describe("Auth Router", func() {
	var (
		router chi.Router
	)

	BeforeEach(func() {
		config.SetForTests(&config.Config{
			JWT: config.JWT{Secret: "secret"},
		})

		router = userhttp.NewAuthRouter(
			dummyUserRepo{},
			noopTx{},
			fixedIDGen{},
			shared.NewSystemClock(),
			notion.NewService(notion.ServiceConfig{}),
		)
	})

	AfterEach(func() {
		config.SetForTests(nil)
	})

	It("sets a state cookie on authorize", func() {
		req := httptest.NewRequest("GET", "/notion/authorize", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)
		Expect(rec.Code).To(Equal(http.StatusOK))

		var resp userhttp.NotionAuthURLResponseDTO
		Expect(json.Unmarshal(rec.Body.Bytes(), &resp)).To(Succeed())
		Expect(resp.State).NotTo(BeEmpty())

		var stateCookie *http.Cookie
		for _, c := range rec.Result().Cookies() {
			if c.Name == "notion_oauth_state" {
				stateCookie = c
			}
		}
		Expect(stateCookie).ToNot(BeNil())
		Expect(stateCookie.Value).To(Equal(resp.State))
		Expect(stateCookie.HttpOnly).To(BeTrue())
	})

	It("rejects callback when state missing or mismatched", func() {
		req := httptest.NewRequest("GET", "/notion/callback?code=abc", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		Expect(rec.Code).To(Equal(http.StatusBadRequest))

		req2 := httptest.NewRequest("GET", "/notion/callback?code=abc&state=wrong", nil)
		req2.AddCookie(&http.Cookie{Name: "notion_oauth_state", Value: "other"})
		rec2 := httptest.NewRecorder()
		router.ServeHTTP(rec2, req2)
		Expect(rec2.Code).To(Equal(http.StatusBadRequest))
		Expect(rec2.Body.String()).To(ContainSubstring("invalid_state"))
	})
})
