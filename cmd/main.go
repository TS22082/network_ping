package main

import (
	"fmt"
	"log"
	"network_testing/internal"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

// Package main implements a network testing application that periodically runs ping tests using ICMP
// and reports results to the Logida API for logging, and analysis.
//
// Before running, ensure you have a .env file with your LOGIDA_API_KEY set.
// The application performs a ping test every 10 minutes and can be gracefully terminated with Ctrl+C.
//
// The ping itself is powered by the pro-bing library - https://github.com/prometheus-community/pro-bing
// The logging dashboard the logs are posted to is powered by Logida - https://logida.fly.dev/
func main() {
	// You will need a .env file with your LOGIDA_API_KEY set.
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load .env: %v: ", err)
	}

	// Retrieve Logida API key from environment variable. Will need to create one at https://logida.fly.dev/
	logidaApiKey := os.Getenv("LOGIDA_API_KEY")
	if logidaApiKey == "" {
		fmt.Println("LOGIDA_API_KEY environment variable not set.")
		return
	}

	// Ticker to trigger periodic tasks every 10 minutes
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	// Channel to listen for interrupt signals (e.g., Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// If you want to run without default setting fill in the PingTestConfig struct here,
	// then remove the .Default() call and use the struct directly in the RunTest call below.
	var config = internal.PingTestConfig{}
	internal.RunTest(config.Default())

	// Main event loop - runs every 10 minutes or until interrupted
	for {
		select {
		case <-ticker.C:
			fmt.Println("Performing periodic task...")
			// if using your own config, replace config.Default() with config
			internal.RunTest(config.Default())
			fmt.Println("Next test schedulled in 10 minutes. Press Ctrl+C to stop.")
		case sig := <-sigChan:
			fmt.Println("Received signal:", sig.String())
			fmt.Println("Shutting down gracefully...")
			return
		}
	}
}
