package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"src/internal/database"
	shared "src/internal/modules/shared/domain"
	"src/internal/modules/users/application"
	"src/internal/modules/users/domain"
	"src/internal/modules/users/infrastructure/postgres"
	"src/internal/pkg/httpx"
)

// NewRouter creates a new HTTP router for the users module
func NewRouter() chi.Router {
	r := chi.NewRouter()

	// Initialize dependencies
	repo := postgres.NewUserRepository(database.GormDB())
	idGen := shared.NewUUIDGenerator()
	clock := shared.NewSystemClock()
	txMgr := shared.NewNoopTransactionManager()

	// Initialize use cases
	createUserUC := application.NewCreateUserUseCase(repo, idGen, clock, txMgr)
	getUserUC := application.NewGetUserUseCase(repo)
	getUserByEmailUC := application.NewGetUserByEmailUseCase(repo)

	// Define routes
	r.Post("/", httpx.EndpointJSON[CreateUserRequestDTO](func(req *http.Request, body CreateUserRequestDTO) (int, any, error) {
		if err := httpx.ValidateTags(body); err != nil {
			return http.StatusUnprocessableEntity, nil, err
		}

		resp, err := createUserUC.Execute(req.Context(), application.CreateUserRequest{
			Email: body.Email,
			Name:  body.Name,
		})
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		dto := toUserResponseDTO(resp.User)
		return http.StatusCreated, dto, nil
	}))

	r.Get("/{userID}", httpx.Endpoint(func(req *http.Request) (int, any, error) {
		userID := chi.URLParam(req, "userID")

		resp, err := getUserUC.Execute(req.Context(), application.GetUserRequest{
			ID: userID,
		})
		if err != nil {
			if err == domain.ErrUserNotFound {
				return http.StatusNotFound, nil, err
			}
			return http.StatusInternalServerError, nil, err
		}

		dto := toUserResponseDTO(resp.User)
		return http.StatusOK, dto, nil
	}))

	r.Get("/by-email/{email}", httpx.Endpoint(func(req *http.Request) (int, any, error) {
		email := chi.URLParam(req, "email")

		resp, err := getUserByEmailUC.Execute(req.Context(), application.GetUserByEmailRequest{
			Email: email,
		})
		if err != nil {
			if err == domain.ErrUserNotFound {
				return http.StatusNotFound, nil, err
			}
			return http.StatusInternalServerError, nil, err
		}

		dto := toUserResponseDTO(resp.User)
		return http.StatusOK, dto, nil
	}))

	return r
}
