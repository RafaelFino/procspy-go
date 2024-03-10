package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"procspy/internal/procspy/domain"
	"procspy/internal/procspy/service"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ValidateRequest is a method to validate a request
//
// It receives a request context and a user service
// It returns a user and an error
// It returns an error if the request is invalid
// It returns an error if the user is not approved
// It returns an error if the user is not authorized
// It returns an error if the key is invalid
// It returns an error if the token is invalid
// It returns an error if the token is expired
// It returns an error if the user data can't be loaded
// It returns an error if the user is not found
// It returns an error if the user key is invalid
// It returns an error if the user key is empty
// It returns an error if the user is empty
// It returns an error if the token is empty
// It returns an error if the user mismatch
// It returns an error if the user key mismatch
// It returns an error if the user data can't be loaded
// It returns an error if the user is not approved
// Request example:
//
//	URL:
//		<url>/:user
//	Header:
//		authorization: <token>
//	Returns:
//		domain.User from token + URL
func ValidateRequest(ctx *gin.Context, userService *service.User, authService *service.Auth) (*domain.User, error) {
	token, err := geToken(ctx)

	if err != nil {
		log.Printf("[Server] Error getting token: %s", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, err
	}

	content, expired, err := authService.Validate(token)

	if err != nil {
		log.Printf("[Server] Error validating request: %s", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, err
	}

	if expired {
		log.Printf("[Server] Token expired")
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "token expired",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, fmt.Errorf("token expired")
	}

	if content == nil {
		log.Printf("[Server] Cannot load user data")
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "cannot load user data",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, fmt.Errorf("cannot load user data")
	}

	user, ok := content["user"]

	if !ok || user == "" {
		log.Printf("[Server] Invalid user")
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid user",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, fmt.Errorf("invalid user")
	}

	key, ok := content["key"]

	if !ok || key == "" {
		log.Printf("[Server] Invalid key")
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid user key",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, fmt.Errorf("invalid key")
	}

	paramUser, err := GetRequestParam(ctx, "user")

	if err != nil {
		log.Printf("[Server] Error getting user: %s", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, err
	}

	if paramUser != user {
		log.Printf("[Server] User mismatch")
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "user mismatch",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, fmt.Errorf("user mismatch")
	}

	userData, err := userService.GetUser(user)

	if err != nil {
		log.Printf("[Server] Error loading user: %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("internal error")
	}

	if !userData.GetApproved() {
		log.Printf("[Server] Unauthorized request - user not approved")
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	if key == "" || userData.GetKey() != key {
		log.Printf("[Server] Unauthorized request - invalid key")
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "unauthorized",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})
		return nil, fmt.Errorf("unauthorized")
	}

	log.Printf("[Server] %s is authorized for this request %s", user, ctx.Request.URL.Path)

	return userData, nil
}

func GetRequestParam(ctx *gin.Context, param string) (string, error) {
	ret := ctx.Param(param)

	if ret == "" {
		log.Printf("[Server] %s is empty", param)
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "param not found"})
		return "", fmt.Errorf("param not found")
	}

	ret = strings.ReplaceAll(ret, "//", "")

	return ret, nil
}

func GetFromBody(ctx *gin.Context, auth *service.Auth, keys []string) (map[string]string, error) {
	body, err := ReadCypherBody(ctx, auth)

	if err != nil {
		log.Printf("[Server] Error reading request body: %s", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":     "invalid request",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		})

		return nil, err
	}

	ret := make(map[string]string, len(keys))

	for _, key := range keys {
		if _, ok := body[key]; !ok || body[key] == "" {
			log.Printf("[Server] %s is empty", key)
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{
				"error":     "invalid request",
				"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
			})
			return nil, fmt.Errorf("invalid request")
		}
		ret[key] = body[key].(string)
	}

	return ret, nil
}
func ReadCypherBody(ctx *gin.Context, auth *service.Auth) (map[string]interface{}, error) {
	ret := make(map[string]interface{}, 0)

	data, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		log.Printf("[Server] Error reading request body: %s", err)
		return ret, err
	}

	jsonData, err := auth.Decypher(string(data))

	if err != nil {
		log.Printf("[Server] Error decyphering request body: %s", err)
		return ret, err
	}

	err = json.Unmarshal([]byte(jsonData), &ret)

	if err != nil {
		log.Printf("[Server] Error parsing request body: %s", err)
	}

	return ret, err
}

func GeNameFromParam(ctx *gin.Context) (string, error) {
	return GetRequestParam(ctx, "name")
}

func geToken(ctx *gin.Context) (string, error) {
	token := ctx.Request.Header.Get("authorization")

	if token == "" {
		log.Printf("[Server] Token is empty")
		return "", fmt.Errorf("token not found")
	}

	return token, nil
}
