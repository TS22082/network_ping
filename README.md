# Network Ping Logger

A lightweight network monitoring tool that periodically runs ICMP ping tests and logs the results to [Logida](https://logida.fly.dev/) for centralized logging and analysis.

## Features

- 🏓 Automated ping tests every 10 minutes
- 📊 Centralized logging via Logida API
- 🛑 Graceful shutdown with Ctrl+C
- ⚙️ Configurable ping parameters (count, interval, target)
- 🔒 Secure API key management via environment variables

## Prerequisites

- Go 1.16 or higher
- A [Logida](https://logida.fly.dev/) account and API key

## Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd network_testing
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the project root:
```bash
LOGIDA_API_KEY=your_api_key_here
```

Get your API key by creating an account at [https://logida.fly.dev/](https://logida.fly.dev/)

## Usage

### Run with default settings
```bash
go run main.go
```

Default configuration:
- **Target:** www.google.com
- **Count:** 100 packets
- **Interval:** 1 second between packets
- **Test frequency:** Every 10 minutes

### Custom configuration

To use custom ping settings, modify the `PingTestConfig` in `main.go`:
```go
config := internal.PingTestConfig{
    Count:    50,                    // Send 50 packets
    Interval: 2 * time.Second,       // 2 seconds between packets
    Target:   "8.8.8.8",             // Ping Google DNS
}
internal.RunTest(config)
```

### Building

Build the executable:
```bash
go build -o network-ping-logger
```

Run the binary:
```bash
./network-ping-logger
```

## Project Structure
```
network_testing/
├── cmd
│   └── main.go                  # Application entry point
├── internal/
│   ├── PingTestConfig.go        # Configuration structure with default method / option
│   ├── SendLog.go               # Sends report to Logida
│   └── RunTest.go               # Ping test implementation
├── .env                         # Environment variables (not committed)
├── go.mod                       # Go module definition
└── README.md
```

## How It Works

1. **Startup:** Loads environment variables and validates Logida API key
2. **Initial Test:** Runs an immediate ping test on startup
3. **Periodic Tests:** Executes ping tests every 10 minutes
4. **Logging:** Sends results to Logida API for centralized monitoring
5. **Shutdown:** Gracefully handles interrupt signals (Ctrl+C, SIGTERM)

## Dependencies

- [pro-bing](https://github.com/prometheus-community/pro-bing) - ICMP ping library
- [godotenv](https://github.com/joho/godotenv) - Environment variable management
- [Logida](https://logida.fly.dev/) - Logging and analytics platform

## Troubleshooting

### "LOGIDA_API_KEY environment variable not set"
- Ensure your `.env` file exists in the project root
- Verify the API key is correctly formatted
- Check file permissions on `.env`

### Tests not running
- Check your internet connection
- Verify the target host is reachable
- Review Logida dashboard for error logs

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[Add your license here]

## Acknowledgments

- [pro-bing](https://github.com/prometheus-community/pro-bing) for the excellent ICMP library
- [Logida](https://logida.fly.dev/) for the logging platform