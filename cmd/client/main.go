package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"procspy/internal/procspy/client"
	"procspy/internal/procspy/config"
	"syscall"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
)

var buildDate string

func main() {
	if len(os.Args) < 2 {
		fmt.Print("Usage: procspy <config_file>\n")
		os.Exit(1)
	}

	configFile := os.Args[1]

	cfg, err := config.ConfigClientFromFile(configFile)
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
	fmt.Print("\nStarting client...\n")

	service := client.NewSpy(cfg)
	go service.Start()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	fmt.Print("\nClient stopped.\n")
}

func initLogger(path string) error {
	if err := os.Mkdir(path, 0755); !os.IsExist(err) {
		fmt.Printf("Error creating directory %s: %s", path, err)
		return err
	}

	//rotate logs every day and store last 30 days
	writer, err := rotatelogs.New(
		fmt.Sprintf("%s/procspy-%%Y%%m%%d.log", path),
		rotatelogs.WithLinkName(fmt.Sprintf("%s/procspy-latest.log", path)),
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
 _____                                                          _____   _   _                  _    
|  __ \                                                        / ____| | | (_)                | |   
| |__) |  _ __    ___     ___   ___   _ __    _   _   ______  | |      | |  _    ___   _ __   | |_  
|  ___/  | '__|  / _ \   / __| / __| | '_ \  | | | | |______| | |      | | | |  / _ \ | '_ \  | __| 
| |      | |    | (_) | | (__  \__ \ | |_) | | |_| |          | |____  | | | | |  __/ | | | | \ |_  
|_|      |_|     \___/   \___| |___/ | .__/   \__, |           \_____| |_| |_|  \___| |_| |_|  \__| 
                                     | |      __/ /                                                 
                                     |_|     |___/                                                  
Build Date: %s
`, buildDate)
}
