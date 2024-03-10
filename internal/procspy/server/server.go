package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	srv *http.Server
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

	gin.ForceConsoleColor()
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

	log.Print("Router started")

	s.srv = &http.Server{
		Addr:    s.config.APIPort,
		Handler: s.router,
	}

	go func() {
		log.Printf("Server running under goroutine, listen and serve on %s", s.config.APIPort)
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}

		log.Print("Server stopped")
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
