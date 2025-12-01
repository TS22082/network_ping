package internal

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func RunTest(cfg PingTestConfig) {
	target := cfg.Target
	count := cfg.Count
	interval := cfg.Interval

	logidaApiKey := os.Getenv("LOGIDA_API_KEY")
	if logidaApiKey == "" {
		fmt.Println("LOGIDA_API_KEY environment variable not set.")
		return
	}

	fmt.Println("Pinging", target, "with", count, "packets at", interval)

	pinger, err := probing.NewPinger(target)
	if err != nil {
		panic(err)
	}

	pinger.Count = count
	pinger.Interval = interval
	pinger.Timeout = time.Duration(count)*interval + time.Second*5
	pinger.SetPrivileged(false)

	pinger.OnSend = func(pkt *probing.Packet) {
		fmt.Printf("Sent packet #%d to %s\n", pkt.Seq, pkt.IPAddr)
	}

	pinger.OnRecv = func(pkt *probing.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n", pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	sixDigitRand := make([]byte, 3)
	_, err = rand.Read(sixDigitRand)
	if err != nil {
		fmt.Println("Error generating random client ID:", err)
		return
	}

	clientIdByTimestamp := fmt.Sprintf("%X", sixDigitRand)

	pinger.OnFinish = func(stats *probing.Statistics) {
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

	if err := pinger.Run(); err != nil {
		panic(err)
	}

	stats := pinger.Statistics()
	fmt.Println(stats.AvgRtt.String())
}
