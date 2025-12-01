package internal

import "time"

// PingTestConfig holds configuration parameters for the RunTest function.
type PingTestConfig struct {
	// Number of ICMP packets to send
	Count int
	// Interval between packets
	Interval time.Duration
	// Target host to ping
	Target string
}

// Default returns a PingTestConfig with standard default values.
// It sets Count to 100 packets, Interval to 1 second, and Target to "www.google.com".
func (u *PingTestConfig) Default() PingTestConfig {
	return PingTestConfig{
		Count:    100,
		Interval: time.Second,
		Target:   "www.google.com",
	}
}
