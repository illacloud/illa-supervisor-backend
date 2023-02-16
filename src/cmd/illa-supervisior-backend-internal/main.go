package main

import (
	"os"

	"github.com/illacloud/illa-supervisior-backend/src/authenticator"
	"github.com/illacloud/illa-supervisior-backend/src/controller"
	"github.com/illacloud/illa-supervisior-backend/src/driver/awss3"
	"github.com/illacloud/illa-supervisior-backend/src/driver/minio"
	"github.com/illacloud/illa-supervisior-backend/src/driver/postgres"
	"github.com/illacloud/illa-supervisior-backend/src/internalrouter"
	"github.com/illacloud/illa-supervisior-backend/src/model"
	"github.com/illacloud/illa-supervisior-backend/src/utils/config"
	"github.com/illacloud/illa-supervisior-backend/src/utils/cors"
	"github.com/illacloud/illa-supervisior-backend/src/utils/logger"
	"github.com/illacloud/illa-supervisior-backend/src/utils/recovery"
	"github.com/illacloud/illa-supervisior-backend/src/utils/tokenvalidator"

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

func initDrive(globalConfig *config.Config, logger *zap.SugaredLogger) *model.Drive {
	if globalConfig.IsAWSDrive() {
		systemAWSConfig := awss3.NewSystemAwsConfigByGlobalConfig(globalConfig)
		teamAWSConfig := awss3.NewTeamAwsConfigByGlobalConfig(globalConfig)
		systemDriveS3Instance := awss3.NewS3Drive(systemAWSConfig)
		teamDriveS3Instance := awss3.NewS3Drive(teamAWSConfig)
		return model.NewDrive(systemDriveS3Instance, teamDriveS3Instance, logger)
	} else if globalConfig.IsMINIODrive() {
		systemMINIOConfig := minio.NewSystemMINIOConfigByGlobalConfig(globalConfig)
		teamMINIOConfig := minio.NewTeamMINIOConfigByGlobalConfig(globalConfig)
		systemDriveS3Instance := minio.NewS3Drive(systemMINIOConfig)
		teamDriveS3Instance := minio.NewS3Drive(teamMINIOConfig)
		return model.NewDrive(systemDriveS3Instance, teamDriveS3Instance, logger)
	}
	logger.Errorw("Error in startup, drive init failed.")
	return nil
}

func initServer() (*Server, error) {
	globalConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	engine := gin.New()
	sugaredLogger := logger.NewSugardLogger()
	// init storage
	postgresConfig, err := postgres.GetPostgresConfig()
	if err != nil {
		return nil, err
	}
	postgresDriver, err := postgres.NewPostgresConnection(postgresConfig, sugaredLogger)
	if err != nil {
		return nil, err
	}
	// init validator
	validator, err := tokenvalidator.NewRequestTokenValidator()
	if err != nil {
		return nil, err
	}
	storage := model.NewStorage(postgresDriver, sugaredLogger)
	drive := initDrive(globalConfig, sugaredLogger)
	c := controller.NewController(storage, drive, validator)
	a := authenticator.NewAuthenticator(storage)
	router := internalrouter.NewRouter(c, a)
	server := NewServer(globalConfig, engine, router, sugaredLogger)
	return server, nil

}

func (server *Server) Start() {
	server.logger.Infow("Starting illa-supervisior-backend-internal.")

	// init
	gin.SetMode(server.config.ServerMode)
	// init cors
	server.engine.Use(gin.CustomRecovery(recovery.CorsHandleRecovery))
	server.engine.Use(cors.Cors())
	server.router.RegisterRouters(server.engine)

	err := server.engine.Run(server.config.ServerHost + ":" + server.config.ServerPort)
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
