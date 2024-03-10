package procspy

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"procspy/internal/procspy/service/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

func LoadFile(path string) (string, error) {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("[util] file not found")
	}

	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("[util] failed to open file")
	}

	defer file.Close()

	buf := make([]byte, stat.Size())
	n, err := file.Read(buf)
	if err != nil {
		return "", fmt.Errorf("[util] failed to read file")
	}

	return string(buf[:n]), nil
}

func WriteFile(path, data string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("[util] failed to create file")
	}

	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return fmt.Errorf("[util] failed to write file")
	}

	return nil
}

func GetRequestParam(c *gin.Context, param string) (string, error) {
	ret := c.Param(param)

	if ret == "" {
		log.Printf("[GetRequestParam] %s is empty", param)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "param not found"})
		return "", fmt.Errorf("param not found")
	}

	ret = strings.ReplaceAll(ret, "//", "")

	return ret, nil
}

func ReadCypherBody(c *gin.Context, auth *auth.Authorization) (map[string]interface{}, error) {
	ret := make(map[string]interface{}, 0)

	data, err := io.ReadAll(c.Request.Body)

	if err != nil {
		log.Printf("[ReadCypherBody] Error reading request body: %s", err)
		return ret, err
	}

	jsonData, err := auth.Decypher(string(data))

	if err != nil {
		log.Printf("[ReadCypherBody] Error decyphering request body: %s", err)
		return ret, err
	}

	err = json.Unmarshal([]byte(jsonData), &ret)

	if err != nil {
		log.Printf("[ReadCypherBody] Error parsing request body: %s", err)
	}

	return ret, err
}

func GetUser(c *gin.Context) (string, error) {
	return GetRequestParam(c, "user")
}

func GeName(c *gin.Context) (string, error) {
	return GetRequestParam(c, "name")
}
