package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	projectsHTTP "src/internal/modules/projects/interfaces/http"
	usersHTTP "src/internal/modules/users/interfaces/http"
	webhooksHTTP "src/internal/modules/webhooks/interfaces/http"
	authmw "src/internal/pkg/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	isLocal := os.Getenv("APP_ENV") == "local"
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	if isLocal {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano}).Level(zerolog.DebugLevel)
	} else {
		logger = logger.Level(zerolog.InfoLevel)
	}

	r.Use(middleware.RequestID)
	r.Use(hlog.NewHandler(logger))
	r.Use(hlog.RequestIDHandler("request_id", "Request-ID"))
	r.Use(hlog.RemoteAddrHandler("remote_ip"))
	r.Use(hlog.UserAgentHandler("user_agent"))
	r.Use(hlog.RefererHandler("referer"))
	r.Use(bodyLogMiddleware(isLocal, 4096))
	r.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", status).
			Int("bytes", size).
			Dur("duration", duration).
			Msg("request")
	}))
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
		r.Mount("/users", usersHTTP.NewRouter(s.userRepo, s.txMgr, s.idGen, s.clock))
		r.Mount("/auth", usersHTTP.NewAuthRouter(s.userRepo, s.txMgr, s.idGen, s.clock, s.notion))

		// Protected routes requiring authentication
		r.Route("/projects", func(r chi.Router) {
			r.Use(authmw.JWTAuthMiddleware)
			r.Mount("/", projectsHTTP.NewRouter(s.projectRepo, s.txMgr, s.idGen, s.clock))
		})

		// Webhook routes with signature validation
		r.Route("/webhooks", func(r chi.Router) {
			r.Mount("/", webhooksHTTP.NewRouter(s.webhookDeps))
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
