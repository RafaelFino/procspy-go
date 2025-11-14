package handlers

import (
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
	start := time.Now()
	user, err := ValidateUser(t.users, ctx)

	if err != nil {
		log.Printf("[handlers.Target.GetTargets] [%s] User validation failed: %v", user, err)
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "user not found",
			"elapsed":   time.Since(start).Milliseconds(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	targets, err := t.service.GetTargets(user)

	if err != nil {
		log.Printf("[handlers.Target.GetTargets] [%s] Failed to retrieve targets from service: %v", user, err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"elapsed":   time.Since(start).Milliseconds(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	matches, err := t.matches.GetMatchesInfo(user)

	if err != nil {
		log.Printf("[handlers.Target.GetTargets] [%s] Failed to retrieve match information: %v", user, err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"elapsed":   time.Since(start).Milliseconds(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	for _, target := range targets.Targets {
		if info, ok := matches[target.Name]; ok {
			target.AddMatchInfo(info)
		}
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"targets":   targets.Targets,
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
