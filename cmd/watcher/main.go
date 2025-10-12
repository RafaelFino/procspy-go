package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"procspy/internal/procspy/config"
	"procspy/internal/procspy/watcher"
	"syscall"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
)

var buildDate string

func main() {
	if len(os.Args) < 2 {
		fmt.Print("Usage: watcher <config_file>\n")
		os.Exit(1)
	}

	configFile := os.Args[1]

	cfg, err := config.WatcherConfigFromFile(configFile)
	if err != nil {
		fmt.Printf("Error loading config file: %s", err)
		os.Exit(1)
	}

	err = initLogger(cfg.LogPath)
	if err != nil {
		fmt.Printf("Error opening log file: %s, using stdout", err)
		log.SetOutput(os.Stdout)
	}

	PrintLogo()
	fmt.Print("\nStarting watcher...\n")

	service := watcher.NewWatcher(cfg)
	go service.Start()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	service.Stop()

	fmt.Print("\nWatcher stopped.\n")
}

func initLogger(path string) error {
	if err := os.Mkdir(path, 0755); !os.IsExist(err) {
		fmt.Printf("Error creating directory %s: %s", path, err)
		return err
	}

	//rotate logs every day and store last 30 days
	writer, err := rotatelogs.New(
		fmt.Sprintf("%s/watcher-%%Y%%m%%d.log", path),
		rotatelogs.WithLinkName(fmt.Sprintf("%s/watcher-latest.log", path)),
		rotatelogs.WithMaxAge(30*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		fmt.Printf("Failed to Initialize Log File %s", err)
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(writer)
	}

	return nil
}

func PrintLogo() {
	fmt.Printf(`
▄▖              ▖  ▖  ▗   ▌     
▙▌▛▘▛▌▛▘▛▘▛▌▌▌▄▖▌▞▖▌▀▌▜▘▛▘▛▌█▌▛▘
▌ ▌ ▙▌▙▖▄▌▙▌▙▌  ▛ ▝▌█▌▐▖▙▖▌▌▙▖▌ 
          ▌ ▄▌                  
Build Date: %s
`, buildDate)
}
