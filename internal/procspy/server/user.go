package server

import (
	"fmt"
	"log"
	"net/http"
	"procspy/internal/procspy"
	"procspy/internal/procspy/domain"
	"procspy/internal/procspy/storage"
	"time"

	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
)

type User struct {
	storage *storage.User
}

func NewUser() *User {
	return &User{}
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

func (u *User) LoadUser(c *gin.Context, content map[string]string) (*domain.User, error) {
	user, ok := content["user"]

	if !ok || user == "" {
		log.Printf("[Server API] Unauthorized request - invalid user")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	userData, err := u.storage.GetUser(user)

	if err != nil {
		log.Printf("[Server API] Error loading user: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("internal error")
	}

	if !userData.GetApproved() {
		log.Printf("[Server API] Unauthorized request - user not approved")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	if key, ok := content["key"]; !ok || key == "" || userData.GetKey() != key {
		log.Printf("[Server API] Unauthorized request - invalid key")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	log.Printf("[Server API] Authorized request for %s", user)

	return userData, nil
}
