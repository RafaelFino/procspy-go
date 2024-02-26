package server

import (
	"fmt"
	"log"
	"net/http"
	"procspy/internal/procspy"
	"procspy/internal/procspy/storage"
	"time"

	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
)

type User struct {
	storage *storage.User
}

func NewUser(dbConn *storage.DbConnection) *User {
	return &User{
		storage: storage.NewUser(dbConn),
	}
}

func (u *User) CreateUser(c *gin.Context) {
	user, err := procspy.GetUser(c)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message":   "user not found",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	key := guuid.New().String()
	err = u.storage.CreateUser(user, key)

	if err != nil {
		log.Printf("[Server API] Error creating user: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message":   "user created",
		"key":       key,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})

	log.Printf("[Server API] User %s created -> %s", user, key)
}
