package procspy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

	wrote, err := file.WriteString(data)
	if err != nil {
		return fmt.Errorf("[util] failed to write file")
	}

	if wrote != len(data) {
		return fmt.Errorf("[util] failed to write all data")
	}

	return nil
}

func DownloadFromURL(url string) (string, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error downloading config: %s", err)
		return "", err
	}
	defer resp.Body.Close()

	buf := new(strings.Builder)
	n, err := io.Copy(buf, resp.Body)

	if n == 0 {
		log.Printf("Error downloading config: %d bytes read: %s", n, err)
		return "", errors.New("no data read")
	}

	if err != nil {
		log.Printf("Error downloading config: %s", err)
		return "", err
	}

	return buf.String(), err
}

func PostToURL(token string, url string, data interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Printf("Error marshalling config: %s", err)
	}

	contentType := "application/json"
	reqBody := []byte(jsonData)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Error posting data: %s", err)
		return "", err
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error posting data: %s", err)
		return "", err
	}
	defer resp.Body.Close()

	buf := new(strings.Builder)
	n, err := io.Copy(buf, resp.Body)

	if n == 0 {
		log.Printf("Error reading data from response: %d bytes read: %s", n, err)
		return "", errors.New("no data read")
	}

	if err != nil {
		log.Printf("Error reading data from response data: %s", err)
		return "", err
	}

	return buf.String(), err
}

func GetLogo() string {
	return `
_____                                                  	
|  __ \                                                	
| |__) |  _ __    ___     ___     ___   _ __    _   _  	
|  ___/  | '__|  / _ \   / __|   / __| | '_ \  | | | | 	
| |      | |    | (_) | | (__    \__ \ | |_) | | |_| | 	
|_|      |_|     \___/   \___|   |___/ | .__/   \__, | 	
                                       | |       __/ | 	
                                       |_|      |___/  	

`
}
