package watcher

import (
	"io"
	"log"
	"net/http"
	"os/exec"
	"procspy/internal/procspy/config"
	"time"
)

type Watcher struct {
	config  *config.Watcher
	enabled bool
}

func NewWatcher(config *config.Watcher) *Watcher {
	ret := &Watcher{config: config}

	return ret
}

func (w *Watcher) Start() {
	log.Printf("[watcher.Start] Watcher service started successfully")

	w.enabled = true

	log.Printf("[watcher.Start] Watcher initialized with configuration:\n%s", w.config.ToJson())

	for w.enabled {
		w.check()
		wait := time.Duration(w.config.Interval) * time.Second
		log.Printf("[watcher.Start] Waiting %s until next health check...", wait)
		time.Sleep(wait)
	}
}

func (w *Watcher) Stop() {
	w.enabled = false
	log.Printf("[watcher.Stop] Watcher service is shutting down...")
}

func executeCommand(command string) (string, error) {
	cmd := exec.Command(command)

	log.Printf("[watcher.executeCommand] Executing command: '%s'", command)
	err := cmd.Run()

	if err != nil {
		log.Printf("[watcher.executeCommand] Failed to execute command '%s': %v", command, err)
	}

	buf, err := cmd.Output()

	if err != nil {
		log.Printf("[watcher.executeCommand] Failed to read command output for '%s': %v", command, err)
	}

	return string(buf), err
}

func (w *Watcher) check() {
	log.Printf("[watcher.check] Performing health check on Procspy service...")

	body, status, err := w.httpGet(w.config.ProcspyURL)

	if err != nil || status != http.StatusOK {
		log.Printf("[watcher.check] Procspy service is down (Status: %d, Error: %v)", status, err)

		if w.config.StartCmd != "" {
			ret, err := executeCommand(w.config.StartCmd)
			if err != nil {
				log.Printf("[watcher.check] Failed to execute start command: %v", err)
			} else {
				log.Printf("[watcher.check] Start command executed successfully. Output: %s", ret)
			}
		} else {
			log.Printf("[watcher.check] No start command configured - unable to restart service")
		}
	} else {
		log.Printf("[watcher.check] Procspy service is healthy (Status: %d, Response: %s)", status, body)
	}
}

func (w *Watcher) httpGet(url string) (string, int, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("[watcher.httpGet] Failed to fetch URL '%s': %v", url, err)
		return "", http.StatusInternalServerError, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[watcher.httpGet] Failed to read response body from '%s': %v", url, err)
		return "", res.StatusCode, err
	}

	return string(body), res.StatusCode, nil
}
