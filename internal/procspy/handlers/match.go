package handlers

import (
	"fmt"
	"log"
	"net/http"
	"procspy/internal/procspy"
	"procspy/internal/procspy/storage"
	"time"

	"github.com/gin-gonic/gin"
)

type Match struct {
	auth    *Auth
	storage *storage.Match
}

func NewMatch(auth *Auth, dbConn *storage.DbConnection) *Match {
	return &Match{
		auth:    auth,
		storage: storage.NewMatch(dbConn),
	}
}

func (m *Match) InsertMatch(c *gin.Context) {
	user, err := m.auth.Validate(c)

	if err != nil {
		log.Printf("[handler.Match] insertMatch -> Error validating request: %s", err)
		return
	}

	if user == nil {
		log.Printf("[handler.Match] insertMatch -> Cannot load user data")
		return
	}

	body, err := procspy.ReadCypherBody(c, m.auth.GetAuth())

	if err != nil {
		log.Printf("[handler.Match] insertMatch -> Error reading request body: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message":   "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	name := body["name"].(string)
	pattern := body["pattern"].(string)
	match := body["match"].(string)
	elapsed := body["elapsed"].(float64)

	err = m.storage.InsertMatch(user.GetName(), name, pattern, match, elapsed)

	if err != nil {
		log.Printf("[handler.Match] insertMatch -> Error inserting match: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	log.Printf("[handler.Match] insertMatch -> Match inserted for %s", user.GetName())

	c.IndentedJSON(http.StatusCreated, gin.H{
		"message":   "match inserted",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}

func (m *Match) GetElapsed(c *gin.Context) {
	user, err := m.auth.Validate(c)

	if err != nil {
		log.Printf("[handler.Match] getElapsed -> Error validating request: %s", err)
		return
	}

	if user == nil {
		log.Printf("[handler.Match] getElapsed -> Cannot load user data")
		return
	}

	matches, err := m.storage.GetElapsed(user.GetName())

	if err != nil {
		log.Printf("[handler.Match] getElapsed -> Error getting elapsed: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	log.Printf("[handler.Match] getElapsed -> %d matches for %s", len(matches), user.GetName())

	c.IndentedJSON(http.StatusOK, gin.H{
		"matches":   matches,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}
