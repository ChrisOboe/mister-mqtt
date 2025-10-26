# WARNING

Most of this repo is AI generated. It was never properly reviewed. Use at you own risk!

# mister-mqtt

An MQTT interface for the MiSTer FPGA platform that provides real-time status updates to HomeAssistant via MQTT discovery.

## Overview

mister-mqtt monitors MiSTer status files and publishes their changes to an MQTT broker using the HomeAssistant MQTT discovery format. This enables seamless integration with HomeAssistant to track:

- Current core name (`/tmp/CORENAME`)
- Active game (`/tmp/ACTIVEGAME`) 
- RBF name (`/tmp/RBFNAME`)

## Features

- **Real-time monitoring**: Uses fsnotify for efficient file system watching
- **HomeAssistant integration**: Automatic MQTT discovery configuration
- **Robust connectivity**: Auto-reconnection and graceful error handling
- **Minimal resource usage**: Event-driven architecture with low overhead
- **Easy configuration**: Simple command-line flag configuration

## Installation

### Prerequisites

- Go 1.25+ (or use Nix development environment)
- Access to an MQTT broker
- MiSTer FPGA system with status files in `/tmp/`

### Using Nix (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd mister-mqtt

# Enter the development environment
nix develop

# Build the application
go build

# Run the application
./mister-mqtt --broker your-mqtt-broker:1883
```

### Manual Build

```bash
# Install dependencies
go mod download

# Build the application
go build

# Run the application
./mister-mqtt --broker your-mqtt-broker:1883
```

## Configuration

Configure the application using command-line flags:

```bash
./mister-mqtt [options]
```

### Available Options

| Flag | Default | Description |
|------|---------|-------------|
| `--broker` | `localhost:1883` | MQTT broker address |
| `--client-id` | `mister-mqtt` | MQTT client ID |
| `--username` | _(empty)_ | MQTT username (optional) |
| `--password` | _(empty)_ | MQTT password (optional) |
| `--topic-prefix` | `homeassistant` | MQTT topic prefix |

**Note**: The device name for HomeAssistant discovery is automatically determined from `/etc/hostname`.

### Example Usage

```bash
# Basic usage with local broker
./mister-mqtt

# With authentication
./mister-mqtt --broker mqtt.example.com:1883 --username myuser --password mypass

# Custom device name and topic prefix
./mister-mqtt --topic-prefix "ha"
```

## HomeAssistant Integration

Once running, mister-mqtt automatically creates three sensors in HomeAssistant:

- **MiSTer Core Name**: Current core being used
- **MiSTer Active Game**: Currently active game
- **MiSTer RBF Name**: Current RBF file name

### MQTT Topics

The application publishes to these MQTT topics (where `{hostname}` is read from `/etc/hostname`):

```
# Discovery topics
homeassistant/sensor/{hostname}/corename/config
homeassistant/sensor/{hostname}/activegame/config
homeassistant/sensor/{hostname}/rbfname/config

# State topics
homeassistant/sensor/{hostname}/corename/state
homeassistant/sensor/{hostname}/activegame/state
homeassistant/sensor/{hostname}/rbfname/state

# Availability topic
homeassistant/sensor/{hostname}/availability
```

## Development

### Development Environment

This project uses Nix flakes for a reproducible development environment:

```bash
# Enter development shell
nix develop

# Available tools:
# - go (Go compiler)
# - gopls (Go language server)
# - go-tools (Go development tools)
# - golangci-lint (Go linter)
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests in verbose mode
go test -v ./...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Vet code
go vet ./...
```

## Architecture

The application consists of three main components:

1. **File Watcher** (`internal/watcher`): Monitors MiSTer status files
2. **MQTT Client** (`internal/mqtt`): Handles MQTT communication and HomeAssistant discovery
3. **Main Application** (`main.go`): Coordinates components and handles configuration

For detailed architecture information, see [docs/architecture.md](docs/architecture.md).

## Troubleshooting

### Common Issues

**Connection refused to MQTT broker**
- Verify broker address and port
- Check network connectivity
- Ensure broker is running and accessible

**Files not being monitored**
- Verify `/tmp/CORENAME`, `/tmp/ACTIVEGAME`, `/tmp/RBFNAME` exist
- Check file permissions
- Ensure MiSTer is writing to these files

**HomeAssistant not discovering sensors**
- Verify MQTT integration is configured in HomeAssistant
- Check topic prefix matches HomeAssistant configuration
- Ensure discovery is enabled in HomeAssistant MQTT integration

### Logging

The application logs important events to stdout:
- MQTT connection status
- File change detections
- Discovery message publications
- Error conditions

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for your changes
4. Ensure tests pass: `go test ./...`
5. Submit a pull request

## License

[Add your license information here]

## Related Links

- [HomeAssistant MQTT Discovery](https://www.home-assistant.io/integrations/mqtt#mqtt-discovery)
- [MiSTer FPGA Project](https://github.com/MiSTer-devel)
- [Eclipse Paho MQTT Go Client](https://github.com/eclipse/paho.mqtt.golang)
- [fsnotify](https://github.com/fsnotify/fsnotify)
