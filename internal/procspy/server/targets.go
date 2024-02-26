package server

import (
	"fmt"
	"log"
	"net/http"
	"procspy/internal/procspy/storage"
	"time"

	"github.com/gin-gonic/gin"
)

type Target struct {
	auth    *Auth
	storage *storage.Target
}

func NewTarget(auth *Auth, dbConn *storage.DbConnection) *Target {
	return &Target{
		auth:          auth,
		targetStorage: storage.NewTarget(dbConn),
	}
}

func (t *Target) getTargets(c *gin.Context) {
	user, err := t.auth.Validate(c)

	if err != nil {
		log.Printf("[Server API] getTargets -> Error validating request: %s", err)
		return
	}

	if user == nil {
		log.Printf("[Server API] getTargets -> Cannot load user data")
		return
	}

	targets, err := t.storage.GetTargets(user.Name)

	if err != nil {
		log.Printf("[Server API] getTargets -> Error getting targets: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	log.Printf("[Server API] getTargets -> %d targets for %s", len(targets), user.GetName())
	c.IndentedJSON(http.StatusOK, gin.H{
		"targets":   targets,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})
}
