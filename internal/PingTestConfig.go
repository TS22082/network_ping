package internal

import "time"

type PingTestConfig struct {
	Count    int
	Interval time.Duration
	Target   string
}

func (u *PingTestConfig) Default() PingTestConfig {
	return PingTestConfig{
		Count:    100,
		Interval: time.Second,
		Target:   "www.google.com",
	}
}
