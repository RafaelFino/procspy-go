package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

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
	s.router.GET("/key", s.getKey)
	s.router.POST("/user/:user", s.createUser)
	s.router.POST("/auth/", s.authenticate)

	s.router.GET("/targets/:user", s.getTargets)
	s.router.POST("/match/:user", s.insertMatch)
	s.router.GET("/match/:user", s.getElapsed)
	s.router.POST("/command/:user/:name", s.logCommand)

	go func() {
		s.router.Run(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
		log.Print("Server stopped")
	}()

	log.Print("Server started")
}

func (s *Server) getTargets(c *gin.Context) {
	content, err := s.auth.V

	if err != nil {
		log.Printf("[Server API] getTargets -> Error validating request: %s", err)
		return
	}

	if user == nil {
		log.Printf("[Server API] getTargets -> Cannot load user data")
		return
	}

	targets, err := s.targetStorage.GetTargets(user.Name)

	if err != nil {
		log.Printf("[Server API] getTargets -> Error getting targets: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	log.Printf("[Server API] getTargets -> %d targets for %s", len(targets), user.GetName())
	c.IndentedJSON(http.StatusOK, gin.H{
		"targets":   targets,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}

func (s *Server) insertMatch(c *gin.Context) {
	user, err := s.validate(c)

	if err != nil {
		log.Printf("[Server API] insertMatch -> Error validating request: %s", err)
		return
	}

	if user == nil {
		log.Printf("[Server API] insertMatch -> Cannot load user data")
		return
	}

	body, err := s.readCypherBody(c)

	if err != nil {
		log.Printf("[Server API] insertMatch -> Error reading request body: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message":   "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	name := body["name"].(string)
	pattern := body["pattern"].(string)
	match := body["match"].(string)
	elapsed := body["elapsed"].(float64)

	err = s.matchStorage.InsertMatch(user.GetName(), name, pattern, match, elapsed)

	if err != nil {
		log.Printf("[Server API] insertMatch -> Error inserting match: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	log.Printf("[Server API] insertMatch -> Match inserted for %s", user.GetName())

	c.IndentedJSON(http.StatusCreated, gin.H{
		"message":   "match inserted",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}

func (s *Server) getElapsed(c *gin.Context) {
	user, err := s.validate(c)

	if err != nil {
		log.Printf("[Server API] getElapsed -> Error validating request: %s", err)
		return
	}

	if user == nil {
		log.Printf("[Server API] getElapsed -> Cannot load user data")
		return
	}

	matches, err := s.matchStorage.GetElapsed(user.GetName())

	if err != nil {
		log.Printf("[Server API] getElapsed -> Error getting elapsed: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	log.Printf("[Server API] getElapsed -> %d matches for %s", len(matches), user.GetName())

	c.IndentedJSON(http.StatusOK, gin.H{
		"matches":   matches,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}

func (s *Server) logCommand(c *gin.Context) {
	user, err := s.validate(c)

	if err != nil {
		log.Printf("[Server API] logCommand -> Error validating request: %s", err)
		return
	}

	if user == nil {
		log.Printf("[Server API] logCommand -> Cannot load user data")
		return
	}

	body, err := s.readCypherBody(c)

	if err != nil {
		log.Printf("[Server API] logCommand -> Error reading request body: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message":   "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	name, err := s.getParam(c, "name")

	if err != nil {
		log.Printf("[Server API] logCommand -> Error getting name: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message":   "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	command := body["command"].(string)
	commandType := body["type"].(string)
	commandReturn := body["return"].(string)

	err = s.CommandStorage.InsertCommand(user.GetName(), name, commandType, command, commandReturn)

	if err != nil {
		log.Printf("[Server API] logCommand -> Error inserting command: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	log.Printf("[Server API] logCommand -> Command inserted for %s::%s", user.GetName(), name)

	c.IndentedJSON(http.StatusCreated, gin.H{
		"message":   "command inserted",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}
