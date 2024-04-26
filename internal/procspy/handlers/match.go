package handlers

import (
	"log"
	"net/http"
	"procspy/internal/procspy/domain"
	"procspy/internal/procspy/service"
	"time"

	"github.com/gin-gonic/gin"
)

type Match struct {
	service *service.Match
	users   *service.Users
}

func NewMatch(matchService *service.Match, usersService *service.Users) *Match {
	return &Match{
		users:   usersService,
		service: matchService,
	}
}

func (m *Match) InsertMatch(ctx *gin.Context) {
	start := time.Now()
	user, err := ValidateUser(m.users, ctx)

	if err != nil {
		log.Printf("[handler.Match] [%s] InsertMatch -> Error validating user: %s", user, err)
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "user not found",
			"elapsed":   time.Since(start).Milliseconds(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	body, err := ctx.GetRawData()

	if err != nil {
		log.Printf("[handler.Match] [%s] InsertMatch -> Error getting raw data: %s", user, err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid json",
			"elapsed":   time.Since(start).Milliseconds(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	match, err := domain.MatchFromJson(string(body))

	if err != nil {
		log.Printf("[handler.Match] [%s] InsertMatch -> Error binding json: %s", user, err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid json",
			"elapsed":   time.Since(start).Milliseconds(),
			"timestamp": time.Now().Format(time.RFC3339),
		})

		return
	}

	err = m.service.InsertMatch(match)

	if err != nil {
		log.Printf("[handler.Match] [%s] InsertMatch -> Error inserting match: %s", user, err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"elapsed":   time.Since(start).Milliseconds(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	log.Printf("[handler.Match] [%s] InsertMatch -> Match Inserted: %s", user, match.ToLog())

	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"message":   "match inserted",
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
