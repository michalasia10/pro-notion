package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"src/internal/config"
	"src/internal/pkg/middleware"
)

func testConfig() *config.Config {
	return &config.Config{
		JWT: struct {
			Secret string
		}{
			Secret: "test-secret-key-for-testing",
		},
	}
}

var _ = Describe("JWT Authentication Middleware", func() {
	var (
		testCfg *config.Config
		userID  uuid.UUID
		token   string
	)

	BeforeEach(func() {
		testCfg = testConfig()
		userID = uuid.New()

		// Temporarily set test config
		config.SetForTests(testCfg)

		// Generate a valid token for testing
		var err error
		token, err = middleware.GenerateJWTToken(userID)
		Expect(err).ToNot(HaveOccurred())
		Expect(token).ToNot(BeEmpty())
	})

	AfterEach(func() {
		// Reset config
		config.SetForTests(nil)
	})

	Describe("JWTAuthMiddleware", func() {
		var (
			nextHandler http.Handler
			req         *http.Request
			rec         *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			nextHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Test handler that verifies user ID was set in context
				actualUserID, err := middleware.GetUserID(r.Context())
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if actualUserID != userID {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				w.WriteHeader(http.StatusOK)
			})
		})

		Context("when no Authorization header is provided", func() {
			It("should return 401 Unauthorized", func() {
				req = httptest.NewRequest("GET", "/api/v1/projects", nil)
				rec = httptest.NewRecorder()

				middleware.JWTAuthMiddleware(nextHandler).ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusUnauthorized))
				Expect(rec.Body.String()).To(ContainSubstring("Missing authorization header"))
			})
		})

		Context("when Authorization header has invalid format", func() {
			It("should return 401 Unauthorized for missing Bearer prefix", func() {
				req = httptest.NewRequest("GET", "/api/v1/projects", nil)
				req.Header.Set("Authorization", token) // Missing "Bearer " prefix
				rec = httptest.NewRecorder()

				middleware.JWTAuthMiddleware(nextHandler).ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusUnauthorized))
				Expect(rec.Body.String()).To(ContainSubstring("Invalid authorization header format"))
			})

			It("should return 401 Unauthorized for empty Bearer token", func() {
				req = httptest.NewRequest("GET", "/api/v1/projects", nil)
				req.Header.Set("Authorization", "Bearer ")
				rec = httptest.NewRecorder()

				middleware.JWTAuthMiddleware(nextHandler).ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusUnauthorized))
				Expect(rec.Body.String()).To(ContainSubstring("Invalid authorization header format"))
			})
		})

		Context("when Authorization header contains invalid token", func() {
			It("should return 401 Unauthorized for malformed JWT", func() {
				req = httptest.NewRequest("GET", "/api/v1/projects", nil)
				req.Header.Set("Authorization", "Bearer invalid.jwt.token")
				rec = httptest.NewRecorder()

				middleware.JWTAuthMiddleware(nextHandler).ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusUnauthorized))
				Expect(rec.Body.String()).To(ContainSubstring("Invalid token"))
			})

			It("should return 401 Unauthorized for JWT signed with wrong key", func() {
				// Create token with different secret
				wrongCfg := &config.Config{
					JWT: struct {
						Secret string
					}{
						Secret: "different-secret",
					},
				}
				config.SetForTests(wrongCfg)
				wrongToken, _ := middleware.GenerateJWTToken(userID)
				config.SetForTests(testCfg) // Reset to test config

				req = httptest.NewRequest("GET", "/api/v1/projects", nil)
				req.Header.Set("Authorization", "Bearer "+wrongToken)
				rec = httptest.NewRecorder()

				middleware.JWTAuthMiddleware(nextHandler).ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusUnauthorized))
				Expect(rec.Body.String()).To(ContainSubstring("Invalid token"))
			})
		})

		Context("when Authorization header contains valid token", func() {
			It("should allow the request and set user ID in context", func() {
				req = httptest.NewRequest("GET", "/api/v1/projects", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec = httptest.NewRecorder()

				middleware.JWTAuthMiddleware(nextHandler).ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("GenerateJWTToken", func() {
		It("should generate a valid JWT token", func() {
			token, err := middleware.GenerateJWTToken(userID)

			Expect(err).ToNot(HaveOccurred())
			Expect(token).ToNot(BeEmpty())

			// Token should be a valid JWT (contain dots)
			parts := strings.Split(token, ".")
			Expect(len(parts)).To(Equal(3))
		})

		It("should generate different tokens for different user IDs", func() {
			userID2 := uuid.New()
			token1, _ := middleware.GenerateJWTToken(userID)
			token2, _ := middleware.GenerateJWTToken(userID2)

			Expect(token1).ToNot(Equal(token2))
		})
	})

	Describe("Context helpers", func() {
		var ctx context.Context

		BeforeEach(func() {
			ctx = context.Background()
		})

		Describe("SetUserID and GetUserID", func() {
			It("should store and retrieve user ID correctly", func() {
				ctxWithUser := middleware.SetUserID(ctx, userID)

				retrievedUserID, err := middleware.GetUserID(ctxWithUser)

				Expect(err).ToNot(HaveOccurred())
				Expect(retrievedUserID).To(Equal(userID))
			})

			It("should return error when user ID is not in context", func() {
				_, err := middleware.GetUserID(ctx)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user ID not found in context"))
			})

			It("should not affect original context", func() {
				ctxWithUser := middleware.SetUserID(ctx, userID)

				// Original context should not have user ID
				_, err := middleware.GetUserID(ctx)
				Expect(err).To(HaveOccurred())

				// New context should have user ID
				retrievedUserID, err := middleware.GetUserID(ctxWithUser)
				Expect(err).ToNot(HaveOccurred())
				Expect(retrievedUserID).To(Equal(userID))
			})
		})
	})
})
