package handlers

import (
	"fmt"
	"log"
	"net/http"
	"procspy/internal/procspy/service"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type User struct {
	auth    *service.Auth
	service *service.User
}

func NewUser(userService *service.User, authService *service.Auth) *User {
	return &User{
		auth:    authService,
		service: userService,
	}
}

// CreateUser is a method to create a user
// It receives a request with a user and key
// It returns an error if the user can't be created
// Body request example:
//
//	{
//		"user": "<user>",
//		"key": "<key>" //optional
//	}
//
// Ok Response Example:
//
//	{
//		"message": "user created",
//		"key": "<key>",
//		"timestamp": "<timestamp>"
//	}
//
// Error Response Example:
//
//	{
//		"error": "internal error",
//		"timestamp": "<timestamp>"
//	}
func (u *User) CreateUser(ctx *gin.Context) {
	body, err := ReadCypherBody(ctx, u.auth)

	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	user := body["user"].(string)

	key, ok := body["key"]

	if !ok {
		key = uuid.NewString()
	}

	err = u.service.CreateUser(user, key.(string))

	if err != nil {
		log.Printf("[handler.User] Error creating user: %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"message":   "user created",
		"key":       key,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})

	log.Printf("[handler.User] User %s created -> %s", user, key)
}
