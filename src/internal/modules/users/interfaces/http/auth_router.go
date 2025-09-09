package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"src/internal/config"
	"src/internal/database"
	shared "src/internal/modules/shared/domain"
	"src/internal/modules/users/application"
	"src/internal/modules/users/domain"
	"src/internal/modules/users/infrastructure/postgres"
	"src/internal/pkg/httpx"
	"src/internal/pkg/notion"
)

// NewAuthRouter creates a new HTTP router for authentication endpoints
func NewAuthRouter() chi.Router {
	r := chi.NewRouter()

	// Initialize dependencies
	cfg := config.Get()
	repo := postgres.NewUserRepository(database.GormDB())
	idGen := shared.NewUUIDGenerator()
	clock := shared.NewSystemClock()
	txMgr := shared.NewNoopTransactionManager()

	// Initialize Notion service
	notionService := notion.NewService(notion.ServiceConfig{
		ClientID:     cfg.Notion.ClientID,
		ClientSecret: cfg.Notion.ClientSecret,
		RedirectURI:  cfg.Notion.RedirectURL,
		APIVersion:   cfg.Notion.APIVersion,
	})

	// Initialize use cases
	getAuthURLUC := application.NewGetAuthorizationURLUseCase(notionService)
	notionOAuthUC := application.NewNotionOAuthUseCase(repo, clock, txMgr, idGen, notionService)

	// Notion OAuth routes
	r.Route("/notion", func(r chi.Router) {
		// GET /api/v1/auth/notion/authorize
		r.Get("/authorize", httpx.Endpoint(func(req *http.Request) (int, any, error) {
			state := req.URL.Query().Get("state")
			if state == "" {
				state = "default" // Generate a proper state in production
			}

			resp, err := getAuthURLUC.Execute(req.Context(), application.GetAuthorizationURLRequest{
				State: state,
			})
			if err != nil {
				return http.StatusInternalServerError, nil, err
			}

			dto := NotionAuthURLResponseDTO{
				AuthorizationURL: resp.AuthorizationURL,
				State:            state,
			}
			return http.StatusOK, dto, nil
		}))

		// GET /api/v1/auth/notion/callback
		r.Get("/callback", httpx.Endpoint(func(req *http.Request) (int, any, error) {
			code := req.URL.Query().Get("code")
			state := req.URL.Query().Get("state")
			errorParam := req.URL.Query().Get("error")

			if errorParam != "" {
				return http.StatusBadRequest, map[string]any{
					"error":   "oauth_error",
					"message": "OAuth authorization failed",
					"details": errorParam,
				}, nil
			}

			if code == "" {
				return http.StatusBadRequest, map[string]any{
					"error":   "missing_code",
					"message": "Authorization code is required",
				}, nil
			}

			resp, err := notionOAuthUC.Execute(req.Context(), application.NotionOAuthRequest{
				Code:  code,
				State: state,
			})
			if err != nil {
				if err == domain.ErrUserNotFound {
					return http.StatusNotFound, nil, err
				}
				return http.StatusInternalServerError, nil, err
			}

			dto := toNotionCallbackResponseDTO(resp.User)
			return http.StatusOK, dto, nil
		}))
	})

	return r
}
