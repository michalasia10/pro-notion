package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	shared "src/internal/modules/shared/domain"
	"src/internal/modules/users/application"
	"src/internal/modules/users/domain"
	"src/internal/pkg/httpx"
	authmw "src/internal/pkg/middleware"
)

// NewRouter creates a new HTTP router for the users module
func NewRouter(repo domain.UserRepository, txMgr shared.TransactionManager, idGen shared.IDGenerator, clock shared.Clock) chi.Router {
	r := chi.NewRouter()

	// Initialize use cases
	createUserUC := application.NewCreateUserUseCase(repo, idGen, clock, txMgr)
	getUserUC := application.NewGetUserUseCase(repo)
	getUserByEmailUC := application.NewGetUserByEmailUseCase(repo)
	getCurrentUserUC := application.NewGetCurrentUserUseCase(repo)

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
			if errors.Is(err, domain.ErrUserAlreadyExists) {
				return 0, nil, httpx.Conflict("user already exists", map[string]any{"email": body.Email})
			}
			return http.StatusInternalServerError, nil, err
		}

		dto := toUserResponseDTO(resp.User)
		return http.StatusCreated, dto, nil
	}))

	r.Group(func(r chi.Router) {
		r.Use(authmw.JWTAuthMiddleware)
		r.Get("/me", httpx.Endpoint(func(req *http.Request) (int, any, error) {
			userID, err := authmw.GetUserID(req.Context())
			if err != nil {
				return http.StatusUnauthorized, nil, err
			}

			resp, err := getCurrentUserUC.Execute(req.Context(), userID.String())
			if err != nil {
				if err == domain.ErrUserNotFound {
					return http.StatusNotFound, nil, err
				}
				return http.StatusInternalServerError, nil, err
			}

			dto := toUserResponseDTO(resp.User)
			return http.StatusOK, dto, nil
		}))
	})

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
