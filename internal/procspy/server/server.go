package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	auth "procspy/internal/procspy/auth"
	config "procspy/internal/procspy/config"
	"procspy/internal/procspy/domain"
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

func (s *Server) getKey(c *gin.Context) {
	key, err := s.auth.GetPubKey()
	if err != nil {
		log.Printf("[Server API] Error getting public key: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"key": key})
}

func (s *Server) getParam(c *gin.Context, param string) (string, error) {
	ret := c.Param(param)

	if ret == "" {
		log.Printf("[Server API] %s is empty", param)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "param not found"})
		return "", fmt.Errorf("param not found")
	}

	ret = strings.ReplaceAll(ret, "//", "")

	return ret, nil
}

func (s *Server) getUserName(c *gin.Context) (string, error) {
	return s.getParam(c, "user")
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

func (s *Server) readCypherBody(c *gin.Context) (map[string]interface{}, error) {
	ret := make(map[string]interface{}, 0)

	data, err := io.ReadAll(c.Request.Body)

	if err != nil {
		log.Printf("[Server API] Error reading request body: %s", err)
		return ret, err
	}

	jsonData, err := s.auth.Decypher(string(data))

	if err != nil {
		log.Printf("[Server API] Error decyphering request body: %s", err)
		return ret, err
	}

	err = json.Unmarshal([]byte(jsonData), &ret)

	if err != nil {
		log.Printf("[Server API] Error parsing request body: %s", err)
	}

	return ret, err
}

func (s *Server) authenticate(c *gin.Context) {
	body, err := s.readCypherBody(c)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message":   "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	requestKey := body["key"].(string)
	requestUser := body["user"].(string)
	requestDate := body["date"].(string)

	if requestKey == "" || requestDate == "" || requestUser == "" {
		log.Printf("[Server API] Error decyphering key")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	authDate, err := time.Parse(time.RFC3339, requestDate)

	if err != nil {
		log.Printf("[Server API] Error parsing date: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid date",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	if authDate.Compare(time.Now().Add(-1*time.Hour)) < 0 {
		log.Printf("[Server API] Unauthorized request - expired token")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	user, err := s.userStorage.GetUser(requestUser)

	if err != nil {
		log.Printf("[Server API] Error getting user: %s", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error":     "user not found",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	if user.GetKey() != requestKey || user.GetApproved() == false {
		log.Printf("[Server API] Unauthorized request")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	content := map[string]string{
		"user": user.Name,
		"key":  user.Key,
	}

	token, err := s.auth.CreateToken(requestUser, content)

	if err != nil {
		log.Printf("[Server API] Error creating token: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	user.SetToken(token)

	c.IndentedJSON(http.StatusOK, gin.H{
		"token":     token,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})

	log.Printf("[Server API] User %s authenticated", requestUser)
}

func (s *Server) validate(c *gin.Context) (*domain.User, error) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		log.Printf("[Server API] Unauthorized request")
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid token",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	content, expired, err := s.auth.Validate(token)

	if err != nil {
		log.Printf("[Server API] Unauthorized request - error: %s", err)
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	if expired {
		log.Printf("[Server API] Unauthorized request - expired token")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "expired token",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	if content == nil {
		log.Printf("[Server API] Unauthorized request - content are nil")
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid token",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	user, ok := content["user"]

	if !ok || user == "" {
		log.Printf("[Server API] Unauthorized request - invalid user")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	userData, err := s.userStorage.GetUser(user)

	if userData.GetApproved() == false {
		log.Printf("[Server API] Unauthorized request - user not approved")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	if key, ok := content["key"]; !ok || key == "" || userData.GetKey() != key {
		log.Printf("[Server API] Unauthorized request - invalid key")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	log.Printf("[Server API] Authorized request for %s", user)

	return userData, nil
}

func (s *Server) getTargets(c *gin.Context) {
	user, err := s.validate(c)

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

	log.Printf("[Server API] getTargets -> %s targets for %s", len(targets), user)
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

	log.Printf("[Server API] insertMatch -> Match inserted for %s", user)

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

	log.Printf("[Server API] getElapsed -> %s matches for %s", len(matches), user.GetName())

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

	log.Printf("[Server API] logCommand -> Command inserted for %s::%s", user, name)

	c.IndentedJSON(http.StatusCreated, gin.H{
		"message":   "command inserted",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}
