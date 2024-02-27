package handlers

import (
	"fmt"
	"log"
	"net/http"
	"procspy/internal/procspy"
	auth "procspy/internal/procspy/auth"
	"procspy/internal/procspy/domain"
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

func (a *Auth) GetPubKey() (c *gin.Context) {
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

func (a *Auth) GetAuth() *auth.Authorization {
	return a.auth
}

func (a *Auth) Authenticate(c *gin.Context) {
	body, err := procspy.ReadCypherBody(c, a.auth)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message":   "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	requestKey := body["key"].(string)
	requestUser := body["user"].(string)
	requestDate := body["date"].(string)

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

	if authDate.Compare(time.Now().Add(-1*time.Hour)) < 0 {
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

func (a *Auth) Validate(c *gin.Context) (*domain.User, error) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		log.Printf("[handler.Auth] Unauthorized request")
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid token",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	content, expired, err := a.auth.Validate(token)

	if err != nil {
		log.Printf("[handler.Auth] Unauthorized request - error: %s", err)
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	if expired {
		log.Printf("[handler.Auth] Unauthorized request - expired token")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "expired token",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	if content == nil {
		log.Printf("[handler.Auth] Unauthorized request - content are nil")
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid token",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	user, ok := content["user"]

	if !ok || user == "" {
		log.Printf("[handler.Auth] Unauthorized request - invalid user")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	userData, err := u.storage.GetUser(user)

	if err != nil {
		log.Printf("[handler.Auth] Error loading user: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("internal error")
	}

	if !userData.GetApproved() {
		log.Printf("[handler.Auth] Unauthorized request - user not approved")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	if key, ok := content["key"]; !ok || key == "" || userData.GetKey() != key {
		log.Printf("[handler.Auth] Unauthorized request - invalid key")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	log.Printf("[handler.Auth] Authorized request for %s", user)

	return userData, nil
}
