package api

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
	auth *auth.Authorization
}

func NewAuth() *Auth {
	return &Auth{
		auth: service.NewAuthorization(),
	}
}

func (a *Auth) GetPubKey() (c *gin.Context) {
	key, err := a.auth.GetPubKey()
	if err != nil {
		log.Printf("[Server API] Error getting public key: %s", err)
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
		log.Printf("[Server API] Error decyphering key")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	authDate, err := time.Parse(time.RFC3339, requestDate)

	if err != nil {
		log.Printf("[Server API] Error parsing date: %s", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid date",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	if authDate.Compare(time.Now().Add(-1*time.Hour)) < 0 {
		log.Printf("[Server API] Unauthorized request - expired token")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	user, err := a.userStorage.GetUser(requestUser)

	if err != nil {
		log.Printf("[Server API] Error getting user: %s", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error":     "user not found",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return
	}

	if user.GetKey() != requestKey || user.GetApproved() == false {
		log.Printf("[Server API] Unauthorized request")
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
		log.Printf("[Server API] Error creating token: %s", err)
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

	log.Printf("[Server API] User %s authenticated", requestUser)
}

func (a *Auth) Validate(c *gin.Context) (*domain.User, error) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		log.Printf("[Server API] Unauthorized request")
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid token",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	content, expired, err := a.auth.Validate(token)

	if err != nil {
		log.Printf("[Server API] Unauthorized request - error: %s", err)
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	if expired {
		log.Printf("[Server API] Unauthorized request - expired token")
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "expired token",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	if content == nil {
		log.Printf("[Server API] Unauthorized request - content are nil")
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid token",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

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
