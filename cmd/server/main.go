package main

import (
	"fmt"
	"log"
	"net/http"

	"cryptoprice-monitor/configs"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("Starting server on port: %d\n", cfg.Server.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
