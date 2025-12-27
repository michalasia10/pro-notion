package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"src/internal/modules/projects/application"
	"src/internal/modules/projects/domain"
	shared "src/internal/modules/shared/domain"
	"src/internal/pkg/httpx"
	"src/internal/pkg/middleware"
)

// NewRouter creates a new HTTP router for the projects module
func NewRouter(
	repo domain.Repository,
	txMgr shared.TransactionManager,
	idGen shared.IDGenerator,
	clock shared.Clock,
) chi.Router {
	r := chi.NewRouter()

	// Initialize use cases
	createProjectUC := application.NewCreateProjectUseCase(repo, idGen, clock, txMgr)

	// Define routes
	r.Post("/", httpx.EndpointJSON[CreateProjectRequestDTO](func(req *http.Request, body CreateProjectRequestDTO) (int, any, error) {
		if err := httpx.ValidateTags(body); err != nil {
			return http.StatusUnprocessableEntity, nil, err
		}

		// Get authenticated user ID from JWT token
		userID, err := middleware.GetUserID(req.Context())
		if err != nil {
			return http.StatusUnauthorized, nil, err
		}

		resp, err := createProjectUC.Execute(req.Context(), application.CreateProjectRequest{
			UserID:              userID,
			NotionDatabaseID:    body.NotionDatabaseID,
			NotionWebhookSecret: body.NotionWebhookSecret,
		})
		if err != nil {
			if errors.Is(err, domain.ErrProjectAlreadyExists) {
				return 0, nil, httpx.Conflict("project already exists", map[string]any{"notionDatabaseId": body.NotionDatabaseID})
			}
			return http.StatusInternalServerError, nil, err
		}

		dto := toProjectResponseDTO(resp.Project)
		return http.StatusCreated, dto, nil
	}))

	r.Get("/", httpx.Endpoint(func(req *http.Request) (int, any, error) {
		// Get authenticated user ID from JWT token
		userID, err := middleware.GetUserID(req.Context())
		if err != nil {
			return http.StatusUnauthorized, nil, err
		}

		projects, err := repo.FindByUserID(req.Context(), userID)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		dto := ProjectsListResponseDTO{
			Projects: toProjectResponseDTOs(projects),
			Count:    len(projects),
		}
		return http.StatusOK, dto, nil
	}))

	return r
}
