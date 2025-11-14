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

	commandHandler     *handlers.Command
	targetHandler      *handlers.Target
	matchHandler       *handlers.Match
	reportHandler      *handlers.Report
	healthcheckHandler *handlers.Healthcheck

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
	log.Printf("[server.initServices] Initializing application services...")
	commandService := service.NewCommand(s.dbConn)
	targetService := service.NewTarget(s.config)
	matchService := service.NewMatch(s.dbConn)
	userService := service.NewUsers(s.config)
	log.Printf("[server.initServices] All services initialized successfully")

	log.Printf("[server.initServices] Initializing HTTP handlers...")
	s.commandHandler = handlers.NewCommand(commandService, userService)
	s.targetHandler = handlers.NewTarget(targetService, userService, matchService)
	s.matchHandler = handlers.NewMatch(matchService, userService)
	s.reportHandler = handlers.NewReport(targetService, userService, matchService, commandService)
	s.healthcheckHandler = handlers.NewHealthcheck()
	log.Printf("[server.initServices] All HTTP handlers initialized successfully")
}

func (s *Server) Start() {
	log.Printf("[server.Start] Starting Procspy server on %s:%d", s.config.APIHost, s.config.APIPort)

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
	s.router.POST("/command/:user", s.commandHandler.InsertCommand)
	s.router.GET("/report/:user", s.reportHandler.GetReport)
	s.router.GET("/healthcheck", s.healthcheckHandler.GetStatus)

	log.Print("[server.Start] HTTP router configured with all endpoints")

	s.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.APIHost, s.config.APIPort),
		Handler: s.router,
	}

	go func() {
		log.Printf("[server.Start] HTTP server listening on %s:%d", s.config.APIHost, s.config.APIPort)
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("[server.Start] Server error: %v", err)
		}

		log.Print("[server.Start] HTTP server stopped")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[server.Start] Received shutdown signal, gracefully shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("[server.Start] Server forced to shutdown due to error:", err)
	}

	log.Println("[server.Start] Server shutdown complete")
}
