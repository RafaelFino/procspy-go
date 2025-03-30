package handlers

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"procspy/internal/procspy/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Report struct {
	service *service.Target
	users   *service.Users
	matches *service.Match
}

func NewReport(targetService *service.Target, usersService *service.Users, matches *service.Match) *Report {
	return &Report{
		service: targetService,
		users:   usersService,
		matches: matches,
	}
}

func (r *Report) GetReport(ctx *gin.Context) {
	start := time.Now()
	user, err := ValidateUser(r.users, ctx)

	if err != nil {
		log.Printf("[handler.Report] [%s] GetReport -> Error validating user: %s", user, err)
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error":     "user not found",
			"elapsed":   time.Since(start).Milliseconds(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	targets, err := r.service.GetTargets(user)

	if err != nil {
		log.Printf("[handler.Report] [%s] GetReport -> Error getting targets: %s", user, err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"elapsed":   time.Since(start).Milliseconds(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	matches, err := r.matches.GetMatchesInfo(user)

	if err != nil {
		log.Printf("[handler.Report] [%s] GetReport -> Error getting matches info: %s", user, err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":     "internal error",
			"elapsed":   time.Since(start).Milliseconds(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	for _, target := range targets.Targets {
		if info, ok := matches[target.Name]; ok {
			target.AddMatchInfo(info)
		}
	}
	htmlContent := ` 		
<html>
<head>
<style>
table {
  font-family: "Noto Sans Mono", monospace;
  border-collapse: collapse;
  width: 100%;
}

td, th {
  border: 1px solid #dddddd;
  text-align: left;
  padding: 8px;
}

tr:nth-child(even) {
  background-color: #dddddd;
}
</style>
</head>
<body>
<h1>Report: ` + user + `</h1>
<table>
<tr><th>Name</th><th>Pattern</th><th>Limit</th><th>Elapsed</th><th>Remaining</th><th>First</th><th>Last</th><th>Weekdays</th><th>Kill</th></tr>`
	for _, target := range targets.Targets {
		htmlContent += "<tr>"
		htmlContent += "<td>" + html.EscapeString(target.Name) + "</td>"
		htmlContent += "<td>" + html.EscapeString(target.Pattern) + "</td>"
		htmlContent += "<td>" + html.EscapeString(FormatInterval(target.Limit)) + "</td>"
		htmlContent += "<td>" + html.EscapeString(FormatInterval(target.Elapsed)) + "</td>"
		htmlContent += "<td>" + html.EscapeString(FormatInterval(target.Remaining)) + "</td>"
		htmlContent += "<td>" + html.EscapeString(target.FirstMatch) + "</td>"
		htmlContent += "<td>" + html.EscapeString(target.LastMatch) + "</td>"
		weekdays := ""
		for k, v := range target.Weekdays {
			weekdays += fmt.Sprintf("%s:%d ", FormatWeedkays(k), int(v))
		}
		htmlContent += "<td>" + html.EscapeString(weekdays) + "</td>"
		htmlContent += "<td>" + strconv.FormatBool(target.Kill) + "</td></tr>"
	}
	htmlContent += "</table></body></html>"

	ctx.Header("Content-Length", strconv.Itoa(len(htmlContent)))
	ctx.Data(http.StatusOK, "text/html", []byte(htmlContent))
}
