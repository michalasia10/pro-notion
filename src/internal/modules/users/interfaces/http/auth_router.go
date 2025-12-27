package http

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	shared "src/internal/modules/shared/domain"
	"src/internal/modules/users/application"
	"src/internal/modules/users/domain"
	"src/internal/pkg/httpx"
	"src/internal/pkg/notion"
)

// NewAuthRouter creates a new HTTP router for authentication endpoints
func NewAuthRouter(
	repo domain.UserRepository,
	txMgr shared.TransactionManager,
	idGen shared.IDGenerator,
	clock shared.Clock,
	notionService *notion.Service,
) chi.Router {
	r := chi.NewRouter()

	// Initialize use cases
	getAuthURLUC := application.NewGetAuthorizationURLUseCase(notionService)
	notionOAuthUC := application.NewNotionOAuthUseCase(repo, clock, txMgr, idGen, notionService)

	// Notion OAuth routes
	r.Route("/notion", func(r chi.Router) {
		// GET /api/v1/auth/notion/authorize
		r.Get("/authorize", func(w http.ResponseWriter, req *http.Request) {
			state, err := generateState()
			if err != nil {
				httpx.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": "failed to generate state"})
				return
			}

			// Persist state in HttpOnly cookie for callback verification
			http.SetCookie(w, &http.Cookie{
				Name:     "notion_oauth_state",
				Value:    state,
				Path:     "/api/v1/auth/notion/callback",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				Expires:  time.Now().Add(10 * time.Minute),
			})

			resp, err := getAuthURLUC.Execute(req.Context(), application.GetAuthorizationURLRequest{
				State: state,
			})
			if err != nil {
				httpx.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
				return
			}

			dto := NotionAuthURLResponseDTO{
				AuthorizationURL: resp.AuthorizationURL,
				State:            state,
			}
			httpx.WriteJSON(w, http.StatusOK, dto)
		})

		// GET /api/v1/auth/notion/callback
		r.Get("/callback", func(w http.ResponseWriter, req *http.Request) {
			code := req.URL.Query().Get("code")
			state := req.URL.Query().Get("state")
			errorParam := req.URL.Query().Get("error")

			if errorParam != "" {
				httpx.WriteJSON(w, http.StatusBadRequest, map[string]any{
					"error":   "oauth_error",
					"message": "OAuth authorization failed",
					"details": errorParam,
				})
				return
			}

			if code == "" {
				httpx.WriteJSON(w, http.StatusBadRequest, map[string]any{
					"error":   "missing_code",
					"message": "Authorization code is required",
				})
				return
			}

			stateCookie, err := req.Cookie("notion_oauth_state")
			if err != nil || state == "" || stateCookie.Value != state {
				httpx.WriteJSON(w, http.StatusBadRequest, map[string]any{
					"error":   "invalid_state",
					"message": "Invalid or missing state parameter",
				})
				return
			}

			resp, err := notionOAuthUC.Execute(req.Context(), application.NotionOAuthRequest{
				Code:  code,
				State: state,
			})
			if err != nil {
				if err == domain.ErrUserNotFound {
					httpx.WriteJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
					return
				}
				httpx.WriteJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
				return
			}

			// Clear state cookie after successful validation
			http.SetCookie(w, expiredStateCookie())

			dto := toNotionCallbackResponseDTO(resp.User, resp.JWTToken)
			httpx.WriteJSON(w, http.StatusOK, dto)
		})
	})

	return r
}

func generateState() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func expiredStateCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "notion_oauth_state",
		Value:    "",
		Path:     "/api/v1/auth/notion/callback",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(-time.Hour),
	}
}
