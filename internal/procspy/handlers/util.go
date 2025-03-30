package handlers

import (
	"errors"
	"log"
	"procspy/internal/procspy/service"

	"github.com/gin-gonic/gin"
)

func ValidateUser(users *service.Users, ctx *gin.Context) (string, error) {
	userName := ctx.Param("user")

	if !users.Exists(userName) {
		log.Printf("[handler.util] ValidateUser -> User %s not found", userName)
		return userName, errors.New("user not found")
	}

	return userName, nil
}
