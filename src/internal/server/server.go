package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"

	"src/internal/config"
	"src/internal/database"
	projectdomain "src/internal/modules/projects/domain"
	userdomain "src/internal/modules/users/domain"
	shared "src/internal/modules/shared/domain"
	"src/internal/pkg/notion"
)

type Server struct {
	port        int
	db          database.Service
	redisClient *redis.Client
	publisher   message.Publisher
	txMgr       shared.TransactionManager
	projectRepo projectdomain.Repository
	userRepo    userdomain.UserRepository
	idGen       shared.IDGenerator
	clock       shared.Clock
	notion      *notion.Service
}

type Dependencies struct {
	DB          database.Service
	Redis       *redis.Client
	Publisher   message.Publisher
	TxMgr       shared.TransactionManager
	ProjectRepo projectdomain.Repository
	UserRepo    userdomain.UserRepository
	IDGen       shared.IDGenerator
	Clock       shared.Clock
	Notion      *notion.Service
}

func NewServer(deps Dependencies) *http.Server {
	cfg := config.Get()

	serverInstance := &Server{
		port:        cfg.Port,
		redisClient: deps.Redis,
		db:          deps.DB,
		publisher:   deps.Publisher,
		txMgr:       deps.TxMgr,
		projectRepo: deps.ProjectRepo,
		userRepo:    deps.UserRepo,
		idGen:       deps.IDGen,
		clock:       deps.Clock,
		notion:      deps.Notion,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverInstance.port),
		Handler:      serverInstance.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
