package server

import (
	"fmt"
	"log"

	"procspy/internal/procspy/config"
	"procspy/internal/procspy/handlers"
	"procspy/internal/procspy/service"
	"procspy/internal/procspy/storage"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine

	config *config.Server

	dbConn *storage.DbConnection

	authHandler    *handlers.Auth
	userHandler    *handlers.User
	commandHandler *handlers.Command
	targetHandler  *handlers.Target
	matchHandler   *handlers.Match
}

func NewServer(config *config.Server) *Server {
	ret := &Server{
		config: config,
		dbConn: storage.NewDbConnection(config),
	}

	ret.initServices()

	return ret
}

func (s *Server) initServices() {
	log.Printf("Creating services...")
	authService := service.NewAuth()
	userService := service.NewUser(s.dbConn)
	commandService := service.NewCommand(s.dbConn)
	targetService := service.NewTarget(s.dbConn)
	matchService := service.NewMatch(s.dbConn)
	log.Printf("Services created")

	log.Printf("Initializing handlers...")
	s.authHandler = handlers.NewAuth(authService, userService)
	s.userHandler = handlers.NewUser(userService, authService)
	s.commandHandler = handlers.NewCommand(commandService, authService, userService)
	s.targetHandler = handlers.NewTarget(targetService, authService, userService)
	s.matchHandler = handlers.NewMatch(matchService, authService, userService)
	log.Printf("Handlers initialized")
}

func (s *Server) Start() {
	log.Printf("Starting server on %s:%d", s.config.Host, s.config.Port)

	gin.DefaultWriter = log.Writer()
	gin.DefaultErrorWriter = log.Writer()

	s.router = gin.Default()
	s.router.GET("/key", s.authHandler.GetPubKey)
	s.router.POST("/user/:user", s.userHandler.CreateUser)
	s.router.POST("/auth/", s.authHandler.Authenticate)

	s.router.GET("/targets/:user", s.targetHandler.GetTargets)
	s.router.POST("/match/:user", s.matchHandler.InsertMatch)
	s.router.GET("/match/:user", s.matchHandler.GetMatches)
	s.router.POST("/command/:user/:name", s.commandHandler.InsertCommand)

	go func() {
		s.router.Run(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
		log.Print("Server stopped")
	}()

	log.Print("Server started")
}

func (s *Server) Stop() {
	s.router = nil
}
