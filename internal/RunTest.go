package internal

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// RunTest performs a ping test based on the provided PingTestConfig
// and reports the results to the Logida API.
func RunTest(cfg PingTestConfig) {
	target := cfg.Target
	count := cfg.Count
	interval := cfg.Interval

	// Retrieve Logida API key from environment variable.
	logidaApiKey := os.Getenv("LOGIDA_API_KEY")
	if logidaApiKey == "" {
		fmt.Println("LOGIDA_API_KEY environment variable not set.")
		return
	}

	// Print test configuration before entering the ping loop
	fmt.Println("Pinging", target, "with", count, "packets at", interval)

	// Create a new pinger instance, the pinger is responsible for sending ICMP packets
	pinger, err := probing.NewPinger(target)
	if err != nil {
		panic(err)
	}

	// Configure the pinger with parameters from PingTestConfig
	pinger.Count = count
	pinger.Interval = interval
	pinger.Timeout = time.Duration(count)*interval + time.Second*5
	pinger.SetPrivileged(false)

	pinger.OnSend = func(pkt *probing.Packet) {
		fmt.Printf("Sent packet #%d to %s\n", pkt.Seq, pkt.IPAddr)
	}

	// pinger.OnRecv (OnFinish, etc.) are callback functions that handle events during the ping test
	// they dont trigger until the pinger.Run() method is called further down.
	pinger.OnRecv = func(pkt *probing.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n", pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	// Generate a unique client ID based on current timestamp, if there is an error present
	// the client ID can be used to track the specific log report the error happened in.
	sixDigitRand := make([]byte, 3)
	_, err = rand.Read(sixDigitRand)
	if err != nil {
		fmt.Println("Error generating random client ID:", err)
		return
	}

	// Convert the random bytes to a hexadecimal string
	clientIdByTimestamp := fmt.Sprintf("%X", sixDigitRand)

	pinger.OnFinish = func(stats *probing.Statistics) {
		// Check for packet loss and report an error if any packets were lost.
		// A log will still happen afterwards with the full statistics. They can be linked
		// together using the client ID.
		if stats.PacketsSent-stats.PacketsRecv > 0 {
			error_report := map[string]interface{}{
				"msg": "track packet loss with client_id",
			}

			error_log := map[string]interface{}{
				"message":  "Error in ping test - packet loss detected",
				"type":     "error",
				"clientId": clientIdByTimestamp,
				"data":     error_report,
			}

			if err := SendLog(error_log); err != nil {
				fmt.Println("Error sending error log report:", err)
			}
		}

		// Prepare the report data to be sent to Logida, Under "data" I can send any key-value pairs
		// I want to include in the log.
		report := map[string]interface{}{
			"Target":        stats.Addr,
			"IP Address":    stats.IPAddr,
			"Total Pings":   stats.PacketsSent,
			"Successful":    stats.PacketsRecv,
			"Failed":        stats.PacketsSent - stats.PacketsRecv,
			"Packet Loss %": stats.PacketLoss,
			"RTT Min":       stats.MinRtt.String(),
			"RTT Avg":       stats.AvgRtt.String(),
			"RTT Max":       stats.MaxRtt.String(),
			"RTT StdDev":    stats.StdDevRtt.String(),
		}

		// Send the log report to Logida. message, type, clientId and data are required fields.
		log := map[string]interface{}{
			"message":  "Ping test completed",
			"type":     "log",
			"clientId": clientIdByTimestamp,
			"data":     report,
		}

		if err := SendLog(log); err != nil {
			fmt.Println("Error sending log report:", err)
		}
	}

	// Start the ping test
	if err := pinger.Run(); err != nil {
		panic(err)
	}

	stats := pinger.Statistics()
	fmt.Println(stats.AvgRtt.String())
}
