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
	service *service.Target
	users   *service.Users
	matches *service.Match
}

func NewTarget(targetService *service.Target, usersService *service.Users, matches *service.Match) *Target {
	return &Target{
		service: targetService,
		users:   usersService,
		matches: matches,
	}
}

func (t *Target) GetTargets(ctx *gin.Context) {
	user, err := ValidateUser(t.users, ctx)

	if err != nil {
		log.Printf("[handler.Target] [%s] GetTargets -> Error validating user: %s", user, err)
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "user not found",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	targets, err := t.service.GetTargets(user)

	if err != nil {
		log.Printf("[handler.Target] [%s] GetTargets -> Error getting targets: %s", user, err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	matches, err := t.matches.GetMatches(user)

	if err != nil {
		log.Printf("[handler.Target] [%s] GetTargets -> Error getting matches: %s", user, err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	for _, target := range targets.Targets {
		if elapsed, ok := matches[target.Name]; ok {
			target.AddElapsed(elapsed)
		}
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"targets":   targets.Targets,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}
