package handlers

import (
	"fmt"
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
	user, err := ValidateUser(m.users, ctx)

	if err != nil {
		log.Printf("[handler.Match] [%s] InsertMatch -> Error validating user: %s", user, err)
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "user not found",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	var match *domain.Match

	if err := ctx.BindJSON(match); err != nil {
		log.Printf("[handler.Match] [%s] InsertMatch -> Error binding json: %s", user, err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid json",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return
	}

	err = m.service.InsertMatch(match)

	if err != nil {
		log.Printf("[handler.Match] [%s] InsertMatch -> Error inserting match: %s", user, err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	log.Printf("[handler.Match] [%s] InsertMatch -> Match Inserted: %s", user, match.ToLog())

	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"message":   "match inserted",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}

func (m *Match) GetMatches(ctx *gin.Context) {
	user, err := ValidateUser(m.users, ctx)

	if err != nil {
		log.Printf("[handler.Match] [%s] GetMatches -> Error validating user: %s", user, err)
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "user not found",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	matches, err := m.service.GetMatches(user)

	if err != nil {
		log.Printf("[handler.Match] [%s] GetMatches -> Error getting matches: %s", user, err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"matches":   matches,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}
