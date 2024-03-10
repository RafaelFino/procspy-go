package handlers

import (
	"fmt"
	"log"
	"net/http"
	"procspy/internal/procspy/domain"
	"procspy/internal/procspy/server"
	"procspy/internal/procspy/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Match struct {
	auth    *service.Auth
	service *service.Match
	user    *service.User
}

func NewMatch(matchService *service.Match, auth *service.Auth, userService *service.User) *Match {
	return &Match{
		auth:    auth,
		user:    userService,
		service: matchService,
	}
}

// GetMatches is a method to get the matches
// It returns the matches
// It returns an error if the matches can't be retrieved
// Ok Response Example:
//
//	{
//		"matches": [
//			{
//				"name": "<name>",
//				"pattern": "<pattern>",
//				"match": "<match>",
//				"elapsed": "<elapsed>"
//			},
//			{
//				"name": "<name>",
//				"pattern": "<pattern>",
//				"match": "<match>",
//				"elapsed": "<elapsed>"
//			}
//		],
//		"timestamp": "<timestamp>"
//	}
//
// Error Response Example:
//
//	{
//		"error": "internal error",
//		"timestamp": "<timestamp>"
//	}
func (m *Match) InsertMatch(ctx *gin.Context) {
	user, err := server.ValidateRequest(ctx, m.user)

	if err != nil {
		log.Printf("[handler.Match] InsertMatch -> Error validating request: %s", err)
		return
	}

	bodyKeys := []string{"name", "pattern", "match", "elapsed"}
	body, err := server.GetFromBody(ctx, m.auth, bodyKeys)

	if err != nil {
		log.Printf("[handler.Match] InsertMatch -> Error reading request body: %s", err)
		return
	}

	name := body["name"]
	pattern := body["pattern"]
	match := body["match"]
	elapsed, err := strconv.ParseFloat(body["elapsed"], 64)

	if err != nil {
		log.Printf("[handler.Match] InsertMatch -> Error parsing elapsed: %s", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	if name != user.GetName() {
		log.Printf("[handler.Match] InsertMatch -> User %s is not allowed to insert match for %s", user.GetName(), name)
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	data := domain.NewMatch(user.GetName(), name, pattern, match, elapsed)

	err = m.service.InsertMatch(data)

	if err != nil {
		log.Printf("[handler.Match] InsertMatch -> Error inserting match: %s (match: %s)", err, data.ToJson())
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	log.Printf("[handler.Match] InsertMatch -> Match inserted for %s", user.GetName())

	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"message":   "match inserted",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}

// GetMatches is a method to get the elapsed matches
// It returns the elapsed matches
// It returns an error if the matches can't be retrieved
// Ok Response Example:
//
//	{
//		"matches": [
//			{
//				"name": "<name>",
//				"pattern": "<pattern>",
//				"match": "<match>",
//				"elapsed": "<elapsed>",
//				"created_at": "<created_at>"
//			},
//			{
//				"name": "<name>",
//				"pattern": "<pattern>",
//				"match": "<match>",
//				"elapsed": "<elapsed>",
//				"created_at": "<created_at>"
//			}
//		],
//		"timestamp": "<timestamp>"
//	}
//
// Error Response Example:
//
//	{
//		"error": "internal error",
//		"timestamp": "<timestamp>"
//	}
func (m *Match) GetMatches(ctx *gin.Context) {
	user, err := server.ValidateRequest(ctx, m.user)

	if err != nil {
		log.Printf("[handler.Match] GetMatches -> Error validating request: %s", err)
		return
	}

	matches, err := m.service.GetElapsed(user.GetName())

	if err != nil {
		log.Printf("[handler.Match] GetMatches -> Error getting elapsed: %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	log.Printf("[handler.Match] GetMatches -> %d matches for %s", len(matches), user.GetName())

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"matches":   matches,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}
