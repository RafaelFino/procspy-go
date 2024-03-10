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

	dbConn         *storage.DbConnection
	authService    *service.Auth
	userService    *service.User
	commandService *service.Command
	targetService  *service.Target
	matchService   *service.Match

	authHandler    *handlers.Auth
	userHandler    *handlers.User
	commandHandler *handlers.Command
	targetHandler  *handlers.Target
	matchHandler   *handlers.Match
}

func NewServer(config *config.Server) *Server {
	return &Server{
		config: config,
	}
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
