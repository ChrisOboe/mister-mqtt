# Architecture Documentation

## Overview

mister-mqtt is a Go application that provides an MQTT interface for the MiSTer FPGA platform. It monitors MiSTer status files and publishes their changes to an MQTT broker using the HomeAssistant MQTT discovery format.

## System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   MiSTer FPGA   │───▶│   mister-mqtt   │───▶│   MQTT Broker   │
│                 │    │                 │    │                 │
│ /tmp/CORENAME   │    │  File Watcher   │    │ HomeAssistant   │
│ /tmp/ACTIVEGAME │    │  MQTT Client    │    │ Integration     │
│ /tmp/RBFNAME    │    │  Discovery Msgs │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Components

### 1. File Watcher (`internal/watcher`)

**Responsibility**: Monitor MiSTer status files for changes using fsnotify.

**Key Features**:
- Watches three files: `/tmp/CORENAME`, `/tmp/ACTIVEGAME`, `/tmp/RBFNAME`
- Creates files if they don't exist
- Provides real-time notifications on file changes
- Handles file system events (write, create)

**Implementation Details**:
- Uses `github.com/fsnotify/fsnotify` for efficient file system monitoring
- Runs in a separate goroutine to avoid blocking the main application
- Graceful error handling for missing or inaccessible files

### 2. MQTT Client (`internal/mqtt`)

**Responsibility**: Handle MQTT communication with the broker and HomeAssistant discovery.

**Key Features**:
- MQTT client configuration and connection management
- HomeAssistant MQTT discovery message generation
- State publishing for MiSTer sensors
- Availability tracking (online/offline status)

**Implementation Details**:
- Uses `github.com/eclipse/paho.mqtt.golang` for MQTT communication
- Supports authentication (username/password)
- Auto-reconnection on connection loss
- Retained messages for device availability

### 3. Main Application (`main.go`)

**Responsibility**: Application entry point, configuration, and coordination.

**Key Features**:
- Command-line flag parsing for configuration
- Component initialization and coordination
- Graceful shutdown handling
- Signal handling for clean exit

## Data Flow

1. **Startup**:
   - Parse command-line flags
   - Initialize MQTT client and connect to broker
   - Publish HomeAssistant discovery messages
   - Initialize file watcher
   - Read initial file values and publish states

2. **Runtime**:
   - File watcher detects changes in MiSTer status files
   - File content is read and processed
   - MQTT client publishes state updates to appropriate topics
   - HomeAssistant receives updates via MQTT

3. **Shutdown**:
   - Signal handler catches interrupt/termination signals
   - File watcher is stopped
   - MQTT client publishes offline status
   - MQTT connection is closed gracefully

## MQTT Topic Structure

### Discovery Topics
```
{topic_prefix}/sensor/{hostname}/corename/config
```

### State Topics
```
{topic_prefix}/sensor/{hostname}/corename/state
```

### Availability Topic
```
{topic_prefix}/sensor/{hostname}/availability
```

Where `{hostname}` is read from `/etc/hostname`.

### Sensor Mappings
- `CORENAME` → `corename` sensor
- `ACTIVEGAME` → `activegame` sensor  
- `RBFNAME` → `rbfname` sensor

## Configuration

The application is configured via command-line flags:

- `--broker`: MQTT broker address (default: localhost:1883)
- `--client-id`: MQTT client ID (default: mister-mqtt)
- `--username`: MQTT username (optional)
- `--password`: MQTT password (optional)
- `--topic-prefix`: MQTT topic prefix (default: homeassistant)

**Note**: The device identifier (node ID) is automatically read from `/etc/hostname` and used in MQTT topic paths and HomeAssistant device identifiers.

## Error Handling

- **File Access Errors**: Logged as warnings, application continues
- **MQTT Connection Errors**: Fatal errors during startup, auto-reconnection during runtime
- **File Watcher Errors**: Logged, watcher continues operating
- **JSON Marshaling Errors**: Logged, discovery publication fails gracefully

## Dependencies

- `github.com/eclipse/paho.mqtt.golang`: MQTT client library
- `github.com/fsnotify/fsnotify`: File system notification library
- Standard Go libraries for configuration, signal handling, and JSON processing

## Testing Strategy

- **Unit Tests**: Cover core functionality of watcher and MQTT components
- **Integration Tests**: Test file watching with temporary files
- **JSON Validation**: Ensure HomeAssistant discovery message format compliance
- **Error Scenarios**: Test handling of missing files and network issues

## Security Considerations

- MQTT credentials can be provided via command-line flags
- File permissions: Application needs read access to `/tmp/CORENAME`, `/tmp/ACTIVEGAME`, `/tmp/RBFNAME`
- Network: Secure MQTT connections depend on broker configuration (TLS support can be added)

## Performance Characteristics

- **Memory Usage**: Minimal, primarily MQTT client buffers and file watchers
- **CPU Usage**: Low, event-driven architecture
- **Network**: Minimal MQTT traffic, only on state changes
- **Disk I/O**: Only when MiSTer updates status files

## Future Enhancements

- TLS/SSL support for secure MQTT connections
- Configuration file support (YAML/JSON)
- Additional MiSTer status file monitoring
- Metrics and health check endpoints
- Docker containerization
- systemd service integration
