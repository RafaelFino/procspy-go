package handlers

import (
	"fmt"
	"log"
	"net/http"
	"procspy/internal/procspy/service"
	"procspy/internal/procspy/storage"
	"time"

	"github.com/gin-gonic/gin"
)

type Command struct {
	auth    *service.Authorization
	storage *storage.Command
}

func NewCommand(auth *Auth, dbConn *storage.DbConnection) *Command {
	return &Command{
		auth:    auth,
		storage: storage.NewCommand(dbConn),
	}
}

func (cm *Command) InsertCommand(c *gin.Context) {
	user, err := cm.auth.Validate(c)

	if err != nil {
		log.Printf("[handler.Command] logCommand -> Error validating request: %s", err)
		return
	}

	if user == nil {
		log.Printf("[handler.Command] logCommand -> Cannot load user data")
		return
	}

	body, err := s.readCypherBody(c)

	if err != nil {
		log.Printf("[handler.Command] logCommand -> Error reading request body: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message":   "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	name, err := s.getParam(c, "name")

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
