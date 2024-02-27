package handlers

import (
	"fmt"
	"log"
	"net/http"
	"procspy/internal/procspy"
	"procspy/internal/procspy/service"
	"time"

	"github.com/gin-gonic/gin"
)

type Command struct {
	auth    *service.Auth
	service *service.Command
}

func NewCommand(authService *service.Auth, commandService *service.Command) *Command {
	return &Command{
		auth:    authService,
		service: commandService,
	}
}

func (c *Command) InsertCommand(ctx *gin.Context) {
	user, err := c.auth.Validate(ctx)

	if err != nil {
		log.Printf("[handler.Command] logCommand -> Error validating request: %s", err)
		return
	}

	if user == nil {
		log.Printf("[handler.Command] logCommand -> Cannot load user data")
		return
	}

	body, err := s.readCypherBody(ctx)

	if err != nil {
		log.Printf("[handler.Command] logCommand -> Error reading request body: %s", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"message":   "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	name, err := procspy.GetName(ctx)

	if err != nil {
		log.Printf("[handler.Command] logCommand -> Error getting name: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message":   "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	command := body["command"].(string)
	commandType := body["type"].(string)
	commandReturn := body["return"].(string)

	err = cm.storage.InsertCommand(user.GetName(), name, commandType, command, commandReturn)

	if err != nil {
		log.Printf("[handler.Command] logCommand -> Error inserting command: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	log.Printf("[handler.Command] logCommand -> Command inserted for %s::%s", user.GetName(), name)

	c.IndentedJSON(http.StatusCreated, gin.H{
		"message":   "command inserted",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}
