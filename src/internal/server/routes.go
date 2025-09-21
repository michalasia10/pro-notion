package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	projectsHTTP "src/internal/modules/projects/interfaces/http"
	usersHTTP "src/internal/modules/users/interfaces/http"
	webhooksHTTP "src/internal/modules/webhooks/interfaces/http"
	authmw "src/internal/pkg/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)

	// API v1 feature routers
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/users", usersHTTP.NewRouter())
		r.Mount("/auth", usersHTTP.NewAuthRouter())

		// Protected routes requiring authentication
		r.Route("/projects", func(r chi.Router) {
			r.Use(authmw.JWTAuthMiddleware)
			r.Mount("/", projectsHTTP.NewRouter())
		})

		// Webhook routes with signature validation
		r.Route("/webhooks", func(r chi.Router) {
			r.Use(authmw.NotionWebhookMiddleware)
			r.Mount("/", webhooksHTTP.NewRouter(s.publisher))
		})
	})

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	health := s.getSystemHealth()
	status := http.StatusOK

	// If any component is unhealthy, return 503
	for _, component := range health {
		if comp, ok := component.(map[string]string); ok {
			if comp["status"] == "down" {
				status = http.StatusServiceUnavailable
				break
			}
		}
	}

	w.WriteHeader(status)
	jsonResp, _ := json.Marshal(health)
	_, _ = w.Write(jsonResp)
}

func (s *Server) getSystemHealth() map[string]interface{} {
	health := make(map[string]interface{})

	// Database health
	dbHealth := s.db.Health()
	health["database"] = dbHealth

	// Redis health
	redisHealth := s.checkRedisHealth()
	health["redis"] = redisHealth

	// Overall system status
	systemHealthy := true
	if dbHealth["status"] == "down" || redisHealth["status"] == "down" {
		systemHealthy = false
	}

	if systemHealthy {
		health["status"] = "healthy"
		health["message"] = "All systems operational"
	} else {
		health["status"] = "unhealthy"
		health["message"] = "One or more systems are down"
	}

	health["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	health["version"] = "1.0.0" // You can make this dynamic

	return health
}

func (s *Server) checkRedisHealth() map[string]string {
	health := make(map[string]string)

	if s.redisClient == nil {
		health["status"] = "down"
		health["error"] = "Redis client not initialized"
		return health
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := s.redisClient.Ping(ctx).Result()
	if err != nil {
		health["status"] = "down"
		health["error"] = fmt.Sprintf("Redis ping failed: %v", err)
		return health
	}

	health["status"] = "up"
	health["message"] = "Redis is healthy"
	return health
}
