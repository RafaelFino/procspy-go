package handlers

import (
	"errors"
	"fmt"
	"log"
	"math"
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

func FormatInterval(seconds float64) string {
	prefix := ""
	if seconds < 0 {
		prefix = "-"
		seconds *= -1
	}

	hours := int(math.Round(seconds / 3600))
	minutes := int((seconds - float64(hours*3600)) / 60)

	return fmt.Sprintf("%s%dh %dm", prefix, hours, minutes)
}

var weekdays = map[int]string{
	0: "Mon",
	1: "Tue",
	2: "Wed",
	3: "Thu",
	4: "Fri",
	5: "Sat",
	6: "Sun",
}

func FormatWeedkays(d int) string {
	d = d % 7
	return weekdays[d]
}
