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
	log.Printf("[Start] Watcher started")

	w.enabled = true

	log.Printf("[Start] Starting with config ->\n%s", w.config.ToJson())

	for w.enabled {
		w.check()
		wait := time.Duration(w.config.Interval) * time.Second
		log.Printf("[Start] Waiting %s until next check...", wait)
		time.Sleep(wait)
	}
}

func (w *Watcher) Stop() {
	w.enabled = false
	log.Printf("[Stop] Stopping...")
}

func executeCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)

	log.Printf("[executeCommand] Executing command: %s", command)
	err := cmd.Run()

	if err != nil {
		log.Printf("[executeCommand] Error executing command: %s -> %s", command, err)
	}

	buf, err := cmd.Output()

	if err != nil {
		log.Printf("[executeCommand] Error to read command output: %s -> %s", command, err)
	}

	return string(buf), err
}

func (w *Watcher) check() {
	log.Printf("[check] Running watcher...")

	body, status, err := w.httpGet(w.config.ProcspyURL)

	if err != nil || status != http.StatusOK {
		log.Printf("[check] Procspy is down! Status: %d, Error: %s", status, err)

		if w.config.StartCmd != "" {
			_, err := executeCommand(w.config.StartCmd)
			if err != nil {
				log.Printf("[check] Error executing start command: %s", err)
			} else {
				log.Printf("[check] Start command executed successfully")
			}
		} else {
			log.Printf("[check] No start command configured")
		}
	} else {
		log.Printf("[check] Procspy is up! Status: %d, Response: %s", status, body)
	}
}

func (w *Watcher) httpGet(url string) (string, int, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("[httpGet] Error getting url: %s", err)
		return "", http.StatusInternalServerError, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[httpGet] Error reading body: %s", err)
		return "", res.StatusCode, err
	}

	log.Printf("[httpGet] %d Response: %s", res.StatusCode, body)

	return string(body), res.StatusCode, nil
}
