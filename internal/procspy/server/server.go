package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	auth "procspy/internal/procspy/auth"
	config "procspy/internal/procspy/config"
	storage "procspy/internal/procspy/storage"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine

	config *config.ServerConfig

	auth *auth.Authorization

	dbConn         *storage.DbConnection
	CommandStorage *storage.Command
	targetStorage  *storage.Target
	userStorage    *storage.User
	matchStorage   *storage.Match
}

func NewServer(config *config.ServerConfig) *Server {
	return &Server{
		config:  config,
		storage: NewServerStorage(config),
		targets: NewServerTargets(config),
		auth:    NewAuthorization(),
	}
}

func (s *Server) Start() {
	log.Printf("Starting server on %s:%d", s.config.Host, s.config.Port)

	gin.DefaultWriter = log.Writer()
	gin.DefaultErrorWriter = log.Writer()

	s.router = gin.Default()
	s.router.GET("/key", s.getKey)
	s.router.GET("/targets/:user", s.getTargets)
	s.router.POST("/user/:user", s.createUser)
	s.router.POST("/auth", s.authenticate)
	s.router.POST("/process/:user", s.insertProcess)
	s.router.POST("/match/:user", s.insertMatch)
	s.router.GET("/elapsed/:user", s.getElapsed)

	go func() {
		s.router.Run(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
		log.Print("Server stopped")
	}()

	log.Print("Server started")
}

func (s *Server) getUserName(c *gin.Context) (string, error) {
	user := c.Param("user")

	if user == "" {
		log.Printf("[Server API] User is empty")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return "", fmt.Errorf("user not found")
	}

	user = strings.ReplaceAll(user, "//", "")

	return user, nil
}

func (s *Server) validateToken(c *gin.Context) (bool, error) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		log.Printf("[Server API] Unauthorized request")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false, nil
	}

	claims, err := s.auth.Validate(token)

	if err != nil {
		log.Printf("[Server API] Unauthorized request - error: %s", err)
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false, err
	}

	if claims == nil {
		log.Printf("[Server API] Unauthorized request - claims are nil")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false, nil
	}

	log.Printf("[Server API] Authorized request for %s", claims.(map[string]interface{})["sub"])

	return true, nil
}

func (s *Server) getKey(c *gin.Context) {
	key, err := s.auth.GetPubKey()
	if err != nil {
		log.Printf("[Server API] Error getting public key: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"key": key})
}

func (s *Server) getTargets(c *gin.Context) {
	auth, err := s.validateToken(c)
	if err != nil {
		log.Printf("Error authorizing: %s", err)
		return
	}

	if !auth {
		return
	}

	user, err := s.getUserName(c)
	if err != nil {
		log.Printf("[Server API] Error getting user: %s", err)
		return
	}

	targets, err := s.targets.GetTargets(user)
	if err != nil {
		log.Printf("[Server API] Error getting targets: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	log.Printf("[Server API] Targets for %s: %v", user, targets)
	c.IndentedJSON(http.StatusOK, gin.H{"targets": targets})
}

func (s *Server) createUser(c *gin.Context) {
	//user string, password string) (error, userKey)
}

func (s *Server) authenticate(c *gin.Context) {
	//userKey string, user string, password string) (string, error)
}

func (s *Server) insertProcess(c *gin.Context) {
	//token string, name string, elapsed float64, pattern string, command string, kill bool) error
}

func (s *Server) insertMatch(c *gin.Context) {
	//token string, name string, pattern string, command string, kill bool
}

func (s *Server) getElapsed(c *gin.Context) {
	//token string) (map[string]float64, error)
}
