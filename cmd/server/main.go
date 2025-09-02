package main

import (
	"cryptoprice-monitor/internal/api"
	"cryptoprice-monitor/internal/fetcher"
	"cryptoprice-monitor/internal/storage"
	"cryptoprice-monitor/internal/util"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"cryptoprice-monitor/configs"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := util.Init(cfg.Logging.EnableFile, cfg.Logging.File); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	util.Info.Println("Starting application")

	cache := storage.NewCache()

	scheduler := fetcher.StartFetcher(cfg, cache)

	router := api.NewRouter(cache)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	util.Info.Printf("Starting server on %s", addr)
	log.Printf("Starting server on %s", addr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(addr, router); err != nil {
			util.Error.Fatalf("Server error: %v", err)
		}
	}()

	<-stop
	util.Info.Println("Shutting down application")
	log.Println("Shutting down application")
	scheduler.Stop()
	util.Info.Println("Scheduler stopped, exiting")
	log.Println("Scheduler stopped, exiting")
}
