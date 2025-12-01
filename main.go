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

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load .env: %v: ", err)
	}

	logidaApiKey := os.Getenv("LOGIDA_API_KEY")
	if logidaApiKey == "" {
		fmt.Println("LOGIDA_API_KEY environment variable not set.")
		return
	}

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var config = internal.PingTestConfig{}
	internal.RunTest(config.Default())
	for {
		select {
		case <-ticker.C:
			fmt.Println("Performing periodic task...")
			internal.RunTest(config.Default())
			fmt.Println("Next test schedulled in 10 minutes. Press Ctrl+C to stop.")
		case sig := <-sigChan:
			fmt.Println("Received signal:", sig.String())
			fmt.Println("Shutting down gracefully...")
			return
		}
	}
}
