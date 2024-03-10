package handlers

import (
	"fmt"
	"log"
	"net/http"
	"procspy/internal/procspy/server"
	"procspy/internal/procspy/service"
	"time"

	"github.com/gin-gonic/gin"
)

type Command struct {
	auth    *service.Auth
	service *service.Command
	user    *service.User
}

func NewCommand(commandService *service.Command, authService *service.Auth, userService *service.User) *Command {
	return &Command{
		auth:    authService,
		service: commandService,
		user:    userService,
	}
}

func (c *Command) InsertCommand(ctx *gin.Context) {
	user, err := server.ValidateRequest(ctx, c.user)

	if err != nil {
		log.Printf("[handler.Command] logCommand -> Error validating request: %s", err)
		return
	}

	bodyKeys := []string{"command", "type", "return"}
	body, err := server.GetFromBody(ctx, c.auth, bodyKeys)

	if err != nil {
		log.Printf("[handler.Command] logCommand -> Error reading request body: %s", err)
		return
	}

	name, err := server.GeNameFromParam(ctx)

	if err != nil {
		log.Printf("[handler.Command] logCommand -> Error getting name: %s", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	command := body["command"]
	commandType := body["type"]
	commandReturn := body["return"]

	err = c.service.InsertCommand(user.GetName(), name, commandType, command, commandReturn)

	if err != nil {
		log.Printf("[handler.Command] logCommand -> Error inserting command: %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	log.Printf("[handler.Command] logCommand -> Command inserted for %s::%s", user.GetName(), name)

	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"message":   "command inserted",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}