package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	auth "procspy/internal/procspy/auth"
	config "procspy/internal/procspy/config"
	storage "procspy/internal/procspy/storage"

	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
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
	s.router.POST("/auth/:user", s.authenticate)

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

func (s *Server) getKey(c *gin.Context) {
	key, err := s.auth.GetPubKey()
	if err != nil {
		log.Printf("[Server API] Error getting public key: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"key": key})
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

func (s *Server) validateToken(c *gin.Context) (bool, map[string]interface{}, error) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		log.Printf("[Server API] Unauthorized request")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false, nil, nil
	}

	claims, err := s.auth.Validate(token)

	if err != nil {
		log.Printf("[Server API] Unauthorized request - error: %s", err)
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false, nil, err
	}

	if claims == nil {
		log.Printf("[Server API] Unauthorized request - claims are nil")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false, nil, nil
	}

	log.Printf("[Server API] Authorized request for %s", claims.(map[string]interface{})["sub"])

	return true, claims.(map[string]interface{}), nil
}

func (s *Server) createUser(c *gin.Context) {
	user, err := s.getUserName(c)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message":   "user not found",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	key := guuid.New().String()
	err = s.userStorage.CreateUser(user, key)

	if err != nil {
		log.Printf("[Server API] Error creating user: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message":   "user created",
		"key":       key,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})

	log.Printf("[Server API] User %s created -> %s", user, key)
}

func (s *Server) authenticate(c *gin.Context) {
	auth, claims, err := s.validateToken(c)

	if err != nil || !auth {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return
	}

	user, err := s.getUserName(c)

}

func (s *Server) getTargets(c *gin.Context) {

}

func (s *Server) insertMatch(c *gin.Context) {

}

func (s *Server) getElapsed(c *gin.Context) {

}

func (s *Server) logCommand(c *gin.Context) {

}
