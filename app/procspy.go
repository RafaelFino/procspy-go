package main

import (
	"fmt"
	"internal/procspy/procspy"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print("Usage: procspy <config_file>\n")
		os.Exit(1)
	}
	configFile := os.Args[1]

	spy := procspy.NewSpy(configFile)
	go spy.Start()
	defer spy.Stop()

	log.Print("Running...")
}
