package main

import (
	"os"

	"github.com/illacloud/illa-supervisor-backend/src/authenticator"
	"github.com/illacloud/illa-supervisor-backend/src/controller"
	"github.com/illacloud/illa-supervisor-backend/src/driver/minio"
	"github.com/illacloud/illa-supervisor-backend/src/driver/postgres"
	"github.com/illacloud/illa-supervisor-backend/src/driver/redis"
	"github.com/illacloud/illa-supervisor-backend/src/internalrouter"
	"github.com/illacloud/illa-supervisor-backend/src/model"
	"github.com/illacloud/illa-supervisor-backend/src/utils/config"
	"github.com/illacloud/illa-supervisor-backend/src/utils/cors"
	"github.com/illacloud/illa-supervisor-backend/src/utils/logger"
	"github.com/illacloud/illa-supervisor-backend/src/utils/recovery"
	"github.com/illacloud/illa-supervisor-backend/src/utils/tokenvalidator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	engine *gin.Engine
	router *internalrouter.Router
	logger *zap.SugaredLogger
	config *config.Config
}

func NewServer(config *config.Config, engine *gin.Engine, router *internalrouter.Router, logger *zap.SugaredLogger) *Server {
	return &Server{
		engine: engine,
		config: config,
		router: router,
		logger: logger,
	}
}

func initStorage(globalConfig *config.Config, logger *zap.SugaredLogger) *model.Storage {
	postgresDriver, err := postgres.NewPostgresConnectionByGlobalConfig(globalConfig, logger)
	if err != nil {
		logger.Errorw("Error in startup, storage init failed.")
	}
	return model.NewStorage(postgresDriver, logger)
}

func initCache(globalConfig *config.Config, logger *zap.SugaredLogger) *model.Cache {
	redisDriver, err := redis.NewRedisConnectionByGlobalConfig(globalConfig, logger)
	if err != nil {
		logger.Errorw("Error in startup, cache init failed.")
	}
	return model.NewCache(redisDriver, logger)

}

func initDrive(globalConfig *config.Config, logger *zap.SugaredLogger) *model.Drive {
	systemMINIOConfig := minio.NewSystemMINIOConfigByGlobalConfig(globalConfig)
	teamMINIOConfig := minio.NewTeamMINIOConfigByGlobalConfig(globalConfig)
	systemDriveS3Instance := minio.NewS3Drive(systemMINIOConfig)
	teamDriveS3Instance := minio.NewS3Drive(teamMINIOConfig)
	return model.NewDrive(systemDriveS3Instance, teamDriveS3Instance, logger)
}

func initServer() (*Server, error) {
	// init
	globalConfig := config.GetInstance()
	engine := gin.New()
	sugaredLogger := logger.NewSugardLogger()

	// init validator
	validator := tokenvalidator.NewRequestTokenValidator()

	// init driver
	storage := initStorage(globalConfig, sugaredLogger)
	cache := initCache(globalConfig, sugaredLogger)
	drive := initDrive(globalConfig, sugaredLogger)

	a := authenticator.NewAuthenticator(storage, cache)
	c := controller.NewController(storage, cache, drive, validator, a)
	router := internalrouter.NewRouter(c, a)
	server := NewServer(globalConfig, engine, router, sugaredLogger)
	return server, nil

}

func (server *Server) Start() {
	server.logger.Infow("Starting illa-supervisor-backend-internal.")

	// init
	gin.SetMode(server.config.ServerMode)
	// init cors
	server.engine.Use(gin.CustomRecovery(recovery.CorsHandleRecovery))
	server.engine.Use(cors.Cors())
	server.router.RegisterRouters(server.engine)

	err := server.engine.Run(server.config.ServerHost + ":" + server.config.InternalServerPort)
	if err != nil {
		server.logger.Errorw("Error in startup", "err", err)
		os.Exit(2)
	}
}

func main() {
	server, err := initServer()

	if err != nil {

	}

	server.Start()
}
