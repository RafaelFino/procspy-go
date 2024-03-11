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

func DownloadFromURL(url string, token string) (map[string]interface{}, int, error) {
	ret := make(map[string]interface{}, 0)
	buf := new(strings.Builder)

	contentType := "application/json"
	reqBody := make([]byte, 0)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("[HTTP.GET] Error requesting data: %s", err)
		return ret, http.StatusInternalServerError, err
	}

	req.Header.Add("Content-Type", contentType)
	if len(token) > 0 {
		req.Header.Add("authorization", token)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[HTTP.GET] Error requesting data: %s", err)
		return ret, http.StatusInternalServerError, err
	}

	defer log.Printf("[HTTP.GET] Response Status: %s from %s -> %s", resp.Status, url, buf.String())
	defer resp.Body.Close()

	n, err := io.Copy(buf, resp.Body)

	if n == 0 {
		log.Printf("[HTTP.GET] Error downloading config: %d bytes read: %s", n, err)
		return ret, resp.StatusCode, errors.New("no data read")
	}

	if err != nil {
		log.Printf("[HTTP.GET] Error downloading config: %s", err)
		return ret, resp.StatusCode, err
	}

	err = json.Unmarshal([]byte(buf.String()), &ret)

	if err != nil {
		log.Printf("[HTTP.GET] Error parsing json: %s", err)
	}

	return ret, resp.StatusCode, err
}

func PostToURL(token string, url string, data interface{}) (map[string]interface{}, int, error) {
	ret := make(map[string]interface{}, 0)
	buf := new(strings.Builder)

	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Printf("[HTTP.POST] Error marshalling config: %s", err)
		return ret, http.StatusInternalServerError, err
	}

	contentType := "application/json"
	reqBody := []byte(jsonData)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("[HTTP.POST] Error posting data: %s", err)
		return ret, http.StatusInternalServerError, err
	}

	req.Header.Add("Content-Type", contentType)
	if len(token) > 0 {
		req.Header.Add("authorization", token)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[HTTP.POST] Error posting data: %s", err)
		return ret, http.StatusInternalServerError, err
	}

	defer log.Printf("[HTTP.POST] Response Status: %s from %s -> %s", resp.Status, url, buf.String())
	defer resp.Body.Close()

	n, err := io.Copy(buf, resp.Body)

	if n == 0 {
		log.Printf("[HTTP.POST] Error reading data from response: %d bytes read: %s", n, err)
		return ret, resp.StatusCode, errors.New("no data read")
	}

	if err != nil {
		log.Printf("[HTTP.POST] Error reading data from response data: %s", err)
		return ret, resp.StatusCode, err
	}

	err = json.Unmarshal([]byte(buf.String()), &ret)

	if err != nil {
		log.Printf("[HTTP.POST] Error parsing json: %s", err)
	}

	return ret, resp.StatusCode, err
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
