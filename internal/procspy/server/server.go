package server

import (
	"context"
	"errors"
	"fmt"
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

	commandHandler *handlers.Command
	targetHandler  *handlers.Target
	matchHandler   *handlers.Match

	srv *http.Server
}

func NewServer(config *config.Server) *Server {
	ret := &Server{
		config: config,
		dbConn: storage.NewDbConnection(config.DBPath),
	}

	ret.initServices()

	return ret
}

func (s *Server) initServices() {
	log.Printf("Creating services...")
	commandService := service.NewCommand(s.dbConn)
	targetService := service.NewTarget(s.config)
	matchService := service.NewMatch(s.dbConn)
	userService := service.NewUsers(s.config)
	log.Printf("Services created")

	log.Printf("Initializing handlers...")
	s.commandHandler = handlers.NewCommand(commandService, userService)
	s.targetHandler = handlers.NewTarget(targetService, userService)
	s.matchHandler = handlers.NewMatch(matchService, userService)
	log.Printf("Handlers initialized")
}

func (s *Server) Start() {
	log.Printf("Starting server on %s:%d", s.config.APIHost, s.config.APIPort)

	gin.ForceConsoleColor()
	gin.DefaultWriter = log.Writer()
	gin.DefaultErrorWriter = log.Writer()
	if s.config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	s.router = gin.Default()
	s.router.GET("/targets/:user", s.targetHandler.GetTargets)
	s.router.POST("/match/:user", s.matchHandler.InsertMatch)
	s.router.GET("/match/:user", s.matchHandler.GetMatches)
	s.router.POST("/command/:user", s.commandHandler.InsertCommand)

	log.Print("Router started")

	s.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.APIHost, s.config.APIPort),
		Handler: s.router,
	}

	go func() {
		log.Printf("Server running under goroutine, listen and serve on %s:%d", s.config.APIHost, s.config.APIPort)
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}

		log.Print("Server stopped")
	}()

	quit := make(chan os.Signal, 1)
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
