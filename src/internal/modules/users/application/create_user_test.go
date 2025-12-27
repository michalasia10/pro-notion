package application_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	shared "src/internal/modules/shared/domain"
	"src/internal/modules/users/application"
	"src/internal/modules/users/domain"
)

type stubUserRepo struct {
	existing    bool
	createErr   error
	createdUser domain.User
}

func (s *stubUserRepo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	s.createdUser = user
	return user, s.createErr
}

func (s *stubUserRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	return domain.User{}, domain.ErrUserNotFound
}

func (s *stubUserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	if s.existing {
		return domain.User{Email: email, ID: uuid.New()}, nil
	}
	return domain.User{}, domain.ErrUserNotFound
}

func (s *stubUserRepo) Update(ctx context.Context, user domain.User) (domain.User, error) {
	return user, nil
}

func (s *stubUserRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *stubUserRepo) List(ctx context.Context, offset, limit int) ([]domain.User, error) {
	return nil, nil
}

type fixedClock struct {
	now time.Time
}

func (f fixedClock) Now() time.Time {
	return f.now
}

type fixedIDGen struct{}

func (fixedIDGen) NewID(prefix string) string { return prefix + "_id" }

var _ = Describe("CreateUserUseCase", func() {
	var (
		repo *stubUserRepo
		uc   *application.CreateUserUseCase
		now  time.Time
	)

	BeforeEach(func() {
		now = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		repo = &stubUserRepo{}
		uc = application.NewCreateUserUseCase(
			repo,
			fixedIDGen{},
			fixedClock{now: now},
			shared.NewNoopTransactionManager(),
		)
	})

	It("returns conflict error when email already exists", func() {
		repo.existing = true

		_, err := uc.Execute(context.Background(), application.CreateUserRequest{
			Email: "duplicate@example.com",
			Name:  "Existing User",
		})

		Expect(err).To(MatchError(domain.ErrUserAlreadyExists))
	})

	It("creates user when email is free", func() {
		resp, err := uc.Execute(context.Background(), application.CreateUserRequest{
			Email: "new@example.com",
			Name:  "New User",
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(resp.User.Email).To(Equal("new@example.com"))
		Expect(resp.User.PublicID).To(Equal("user_id"))
		Expect(resp.User.CreatedAt).To(Equal(now))
		Expect(repo.createdUser.CreatedAt).To(Equal(now))
	})
})
