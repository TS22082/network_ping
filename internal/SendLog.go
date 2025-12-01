package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

func SendLog(logData map[string]interface{}) error {
	logidaApiKey := os.Getenv("LOGIDA_API_KEY")
	if logidaApiKey == "" {
		println("LOGIDA_API_KEY environment variable not set.")
		return errors.New("LOGIDA_API_KEY not set")
	}

	logAsJSON, err := json.MarshalIndent(logData, "", "  ")
	if err != nil {
		fmt.Println("Error generating JSON log:", err)
		return err
	}

	logURL := "https://logida.fly.dev/api/log"
	body := bytes.NewBuffer(logAsJSON)

	req, err := http.NewRequest("POST", logURL, body)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", logidaApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending log report:", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send log report, status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	return nil
}
