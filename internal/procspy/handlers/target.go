package handlers

import (
	"fmt"
	"log"
	"net/http"
	"procspy/internal/procspy/service"
	"time"

	"github.com/gin-gonic/gin"
)

type Target struct {
	auth    *service.Auth
	service *service.Target
	user    *service.User
}

func NewTarget(targetService *service.Target, authService *service.Auth, userService *service.User) *Target {
	return &Target{
		auth:    authService,
		service: targetService,
		user:    userService,
	}
}

// GetTargets is a method to get the targets
// It returns the targets
// It returns an error if the targets can't be retrieved
// Ok Response Example:
//
//	{
//		"targets": [
//			{
//				"user": "<user>",
//				"name": "<name>",
//				"pattern": "<pattern>",
//				"elapsed": "<elapsed>",
//				"limit": "<limit>",
//				"kill": "<kill>",
//				"so_source": "<so_source>",
//				"check_cmd": "<check_cmd>",
//				"warn_cmd": "<warn_cmd>",
//				"elapsed_cmd": "<elapsed_cmd>",
//			}
//		],
//		"timestamp": "<timestamp>"
//	}
func (t *Target) GetTargets(ctx *gin.Context) {
	user, err := ValidateRequest(ctx, t.user, t.auth)

	if err != nil {
		log.Printf("[handler.Match] GetTargets -> Error validating request: %s", err)
		return
	}

	targets, err := t.service.GetTargets(user.GetName())

	if err != nil {
		log.Printf("[handler.Target] GetTargets -> Error getting targets: %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	log.Printf("[handler.Target] GetTargets -> %d targets for %s", len(targets), user.GetName())
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"targets":   targets,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}
