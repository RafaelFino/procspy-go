package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Healthcheck struct {
	startTime time.Time
}

func NewHealthcheck() *Healthcheck {
	return &Healthcheck{
		startTime: time.Now(),
	}
}

func (h *Healthcheck) GetStatus(ctx *gin.Context) {
	log.Printf("[handlers.Healthcheck.GetStatus] Health check passed - server is running")
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"elapsed":   time.Since(h.startTime).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
