package server

import (
	"fmt"
	"log"

	"procspy/internal/procspy/config"
	"procspy/internal/procspy/storage"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine

	config *config.ServerConfig

	dbConn         *storage.DbConnection
	authService    *Auth
	userService    *User
	commandService *Command
	targetService  *Target
	matchService   *Match
}

func NewServer(config *config.ServerConfig) *Server {
	return &Server{
		config: config,
	}
}

func (s *Server) Start() {
	log.Printf("Starting server on %s:%d", s.config.Host, s.config.Port)

	gin.DefaultWriter = log.Writer()
	gin.DefaultErrorWriter = log.Writer()

	s.router = gin.Default()
	s.router.GET("/key", s.authService.GetPubKey)
	s.router.POST("/user/:user", s.userService.CreateUser)
	s.router.POST("/auth/", s.authService.Authenticate)

	s.router.GET("/targets/:user", s.targetService.GetTargets)
	s.router.POST("/match/:user", s.matchService.InsertMatch)
	s.router.GET("/match/:user", s.matchService.GetMatches)
	s.router.POST("/command/:user/:name", s.commandService.InsertCommand)

	go func() {
		s.router.Run(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
		log.Print("Server stopped")
	}()

	log.Print("Server started")
}
