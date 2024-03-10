package handlers

import (
	"fmt"
	"log"
	"net/http"
	"procspy/internal/procspy/server"
	"procspy/internal/procspy/service"
	"time"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	auth *service.Auth
	user *service.User
}

func NewAuth(authService *service.Auth, userService *service.User) *Auth {
	return &Auth{
		auth: authService,
		user: userService,
	}
}

// GetPubKey is a method to get the public key
// It returns the public key
// It returns an error if the public key can't be retrieved
// Ok Response Example:
//
//	{
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
func (a *Auth) GetPubKey(c *gin.Context) {
	key, err := a.auth.GetPubKey()
	if err != nil {
		log.Printf("[handler.Auth] Error getting public key: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"key":       key,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})

	return
}

// Authenticate is a method to authenticate a user
//
// It receives a request with a key, user and date
// It returns a token if the user is authenticated
// It returns an error if the request is invalid
// It returns an error if the user is not found
// It returns an error if the user is not approved
// It returns an error if the key is invalid
// It returns an error if the token can't be created
// It returns an error if the date is invalid
// It returns an error if the token is expired
// It returns an error if the request is unauthorized
//
// Expected body Example:
//
//	{
//		"key": "key",
//		"user": "user",
//		"date": "2021-01-01T00:00:00Z"
//	}
//
// Ok response Example:
//
//	{
//		"token": "<token>",
//		"timestamp": "<timestamp>"
//	}
//
// Error response Example:
//
//	{
//		"error": "unauthorized",
//		"timestamp": "<timestamp>"
//	}
func (a *Auth) Authenticate(c *gin.Context) {
	bodyKeys := []string{"key", "user", "date"}
	body, err := server.GetFromBody(c, a.auth, bodyKeys)

	requestKey := body["key"]
	requestUser := body["user"]
	requestDate := body["date"]

	if requestKey == "" || requestDate == "" || requestUser == "" {
		log.Printf("[handler.Auth] Error decyphering key")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	authDate, err := time.Parse(time.RFC3339, requestDate)

	if err != nil {
		log.Printf("[handler.Auth] Error parsing date: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid date",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	if authDate.Compare(time.Now().Add(-24*time.Hour)) < 0 {
		log.Printf("[handler.Auth] Unauthorized request - expired token")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	user, err := a.user.GetUser(requestUser)

	if err != nil {
		log.Printf("[handler.Auth] Error getting user: %s", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error":     "user not found",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	if user.GetKey() != requestKey || user.GetApproved() == false {
		log.Printf("[handler.Auth] Unauthorized request")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	content := map[string]string{
		"user": user.Name,
		"key":  user.Key,
	}

	token, err := a.auth.CreateToken(requestUser, content)

	if err != nil {
		log.Printf("[handler.Auth] Error creating token: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	user.SetToken(token)

	c.IndentedJSON(http.StatusOK, gin.H{
		"token":     token,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	})

	log.Printf("[handler.Auth] User %s authenticated", requestUser)
}
