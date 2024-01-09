package main

import (
	"fmt"
	"os"
	"procspy/internal/procspy"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print("Usage: procspy <config_file>\n")
		os.Exit(1)
	}

	fmt.Print("Starting...\n")
	configFile := os.Args[1]

	spy := procspy.NewSpy(configFile)
	go spy.Start()
	defer spy.Stop()

	fmt.Print("Press enter to stop...\n")
}
