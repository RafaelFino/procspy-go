package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"procspy/internal/procspy"
	"procspy/internal/procspy/config"
	"procspy/internal/procspy/server"
	"syscall"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print("Usage: procspy-server <config_file>\n")
		os.Exit(1)
	}

	configFile := os.Args[1]

	jsonConfig, err := procspy.LoadFile(configFile)

	if err != nil {
		fmt.Printf("Error loading config file: %s", err)
		os.Exit(1)
	}

	cfg, err := config.ServerFromJson(jsonConfig)
	if err != nil {
		fmt.Printf("Error loading config file: %s", err)
		os.Exit(1)
	}

	err = initLogger(cfg.LogPath)
	if err != nil {
		fmt.Printf("Error opening log file: %s, using stdout", err)
		log.SetOutput(os.Stdout)
	}

	fmt.Printf("%s\nStarting...", procspy.GetLogo())

	service := server.NewServer(cfg)
	go service.Start()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	log.Print("Stopping...\n")
}

func initLogger(path string) error {
	if err := os.Mkdir(path, 0755); !os.IsExist(err) {
		fmt.Printf("Error creating directory %s: %s", path, err)
		return err
	}

	writer, err := rotatelogs.New(
		fmt.Sprintf("%s/%s.log", path, "%Y%m%d"),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithRotationCount(30), //30 days
	)
	if err != nil {
		log.Fatalf("Failed to Initialize Log File %s", err)
	}
	log.SetOutput(writer)

	return nil
}
