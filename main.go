package main

import (
	"fmt"
	"log"

	"gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	err = cfg.SetUser("liam") // change "liam" to your name if needed
	if err != nil {
		log.Fatalf("failed to set user: %v", err)
	}

	updatedCfg, err := config.Read()
	if err != nil {
		log.Fatalf("failed to read updated config: %v", err)
	}

	fmt.Printf("Config:\nDB URL: %s\nCurrent User: %s\n", updatedCfg.DBUrl, updatedCfg.CurrentUserName)
}

