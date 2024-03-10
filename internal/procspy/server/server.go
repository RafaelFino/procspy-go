package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"procspy/internal/procspy/config"
	"procspy/internal/procspy/domain"
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

// ValidateRequest is a method to validate a request
//
// It receives a request context and a user service
// It returns a user and an error
// It returns an error if the request is invalid
func ValidateRequest(ctx *gin.Context, userService *service.User) (*domain.User, error) {
	token, err := geToken(ctx)

	if err != nil {
		log.Printf("[Server] Error getting token: %s", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, err
	}

	content, expired, err := s.authService.Validate(token)

	if err != nil {
		log.Printf("[Server] Error validating request: %s", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, err
	}

	if expired {
		log.Printf("[Server] Token expired")
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "token expired",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, fmt.Errorf("token expired")
	}

	if content == nil {
		log.Printf("[Server] Cannot load user data")
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "cannot load user data",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, fmt.Errorf("cannot load user data")
	}

	user, ok := content["user"]

	if !ok || user == "" {
		log.Printf("[Server] Invalid user")
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid user",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, fmt.Errorf("invalid user")
	}

	key, ok := content["key"]

	if !ok || key == "" {
		log.Printf("[Server] Invalid key")
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid user key",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, fmt.Errorf("invalid key")
	}

	paramUser, err := GetRequestParam(ctx, "user")

	if err != nil {
		log.Printf("[Server] Error getting user: %s", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, err
	}

	if paramUser != user {
		log.Printf("[Server] User mismatch")
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "user mismatch",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, fmt.Errorf("user mismatch")
	}

	userData, err := userService.GetUser(user)

	if err != nil {
		log.Printf("[Server] Error loading user: %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("internal error")
	}

	if !userData.GetApproved() {
		log.Printf("[Server] Unauthorized request - user not approved")
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	if key == "" || userData.GetKey() != key {
		log.Printf("[Server] Unauthorized request - invalid key")
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	log.Printf("[Server] %s is authorized for this request %s", user, ctx.Request.URL.Path)

	return userData, nil
}

func GetRequestParam(c *gin.Context, param string) (string, error) {
	ret := c.Param(param)

	if ret == "" {
		log.Printf("[Server] %s is empty", param)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "param not found"})
		return "", fmt.Errorf("param not found")
	}

	ret = strings.ReplaceAll(ret, "//", "")

	return ret, nil
}

func GetFromBody(c *gin.Context, auth *service.Auth, keys []string) (map[string]string, error) {
	body, err := ReadCypherBody(c, auth)

	if err != nil {
		log.Printf("[Server] Error reading request body: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, err
	}

	ret := make(map[string]string, len(keys))

	for _, key := range keys {
		if _, ok := body[key]; !ok || body[key] == "" {
			log.Printf("[Server] %s is empty", key)
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"error":     "invalid request",
				"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
			})
			return nil, fmt.Errorf("invalid request")
		}
		ret[key] = body[key].(string)
	}

	return ret, nil
}
func ReadCypherBody(c *gin.Context, auth *service.Auth) (map[string]interface{}, error) {
	ret := make(map[string]interface{}, 0)

	data, err := io.ReadAll(c.Request.Body)

	if err != nil {
		log.Printf("[Server] Error reading request body: %s", err)
		return ret, err
	}

	jsonData, err := auth.Decypher(string(data))

	if err != nil {
		log.Printf("[Server] Error decyphering request body: %s", err)
		return ret, err
	}

	err = json.Unmarshal([]byte(jsonData), &ret)

	if err != nil {
		log.Printf("[Server] Error parsing request body: %s", err)
	}

	return ret, err
}

func GeNameFromParam(c *gin.Context) (string, error) {
	return GetRequestParam(c, "name")
}

func geToken(c *gin.Context) (string, error) {
	token := c.Request.Header.Get("authorization")

	if token == "" {
		log.Printf("[Server] Token is empty")
		return "", fmt.Errorf("token not found")
	}

	return token, nil
}
