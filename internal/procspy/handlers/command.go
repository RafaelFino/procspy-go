package handlers

import (
	"log"
	"net/http"
	"procspy/internal/procspy/domain"
	"procspy/internal/procspy/service"
	"time"

	"github.com/gin-gonic/gin"
)

type Command struct {
	service *service.Command
	users   *service.Users
}

func NewCommand(commandService *service.Command, usersService *service.Users) *Command {
	return &Command{
		service: commandService,
		users:   usersService,
	}
}

func (c *Command) InsertCommand(ctx *gin.Context) {
	start := time.Now()
	user, err := ValidateUser(c.users, ctx)

	if err != nil {
		log.Printf("[handler.Command] [%s] InsertCommand -> Error validating user: %s", user, err)
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "user not found",
			"elapsed":   time.Since(start).Milliseconds(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	body, err := ctx.GetRawData()

	if err != nil {
		log.Printf("[handler.Command] [%s] InsertCommand -> Error getting raw data: %s", user, err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid json",
			"elapsed":   time.Since(start).Milliseconds(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	cmd, err := domain.CommandFromJson(string(body))

	if err != nil {
		log.Printf("[handler.Command] [%s] InsertCommand -> Error binding json: %s", user, err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid json",
			"elapsed":   time.Since(start).Milliseconds(),
			"timestamp": time.Now().Format(time.RFC3339),
		})

		return
	}

	err = c.service.InsertCommand(cmd)

	if err != nil {
		log.Printf("[handler.Command] [%s] InsertCommand -> Error inserting command: %s", user, err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"elapsed":   time.Since(start).Milliseconds(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	log.Printf("[handler.Command] [%s] InsertCommand -> Command inserted: %s", user, cmd.ToLog())

	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"message":   "command inserted",
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
