package config

import (
	"encoding/json"
	"log"
	"os"
)

type Watcher struct {
	Interval   int    `json:"interval"`
	LogPath    string `json:"log_path"`
	ProcspyURL string `json:"procspy_url"`
	StartCmd   string `json:"start_cmd"`
}

func NewWatcher() *Watcher {
	return &Watcher{
		Interval:   10,
		LogPath:    "logs",
		ProcspyURL: "http://localhost:8888",
		StartCmd:   "",
	}
}

func (w *Watcher) SetDefaults() {
	if w.Interval < 10 {
		w.Interval = 10
	}

	if w.LogPath == "" {
		w.LogPath = "logs"
	}

	if w.ProcspyURL == "" {
		w.ProcspyURL = "http://localhost:8888"
	}
}

func (w *Watcher) ToJson() string {
	ret, err := json.MarshalIndent(w, "", "\t")
	if err != nil {
		log.Printf("[config.Watcher.ToJson] Failed to marshal watcher configuration to JSON: %v", err)
	}

	return string(ret)
}

func WatcherConfigFromJson(jsonString string) (*Watcher, error) {
	ret := &Watcher{}
	err := json.Unmarshal([]byte(jsonString), ret)
	if err != nil {
		log.Printf("[config.WatcherConfigFromJson] Failed to unmarshal watcher configuration: %v", err)
		return nil, err
	}

	ret.SetDefaults()

	log.Printf("[config.WatcherConfigFromJson] Watcher configuration loaded successfully: %s", ret.ToJson())

	return ret, nil
}

func WatcherConfigFromFile(path string) (*Watcher, error) {
	byteValue, err := os.ReadFile(path)
	if err != nil {
		log.Printf("[config.WatcherConfigFromFile] Failed to read configuration file '%s': %v", path, err)
		return nil, err
	}

	return WatcherConfigFromJson(string(byteValue))
}
