
# EnSync CLI Tool

A command-line interface for interacting with the EnSync config management service.

---

## Table of Contents
1. [Installation](#installation)
2. [Configuration](#configuration)
3. [Usage](#usage)
    - [Event Management](#event-management)
    - [Access Key Management](#access-key-management)
    - [General Options](#general-options)
4. [Common Flags](#common-flags)
5. [Error Handling](#error-handling)
6. [Development](#development)

---

## Installation

### 1. From Source
```bash
# Clone the repository
git clone https://github.com/rossi1/ensync-cli
cd ensync-cli

# Build the binary
make build

# The binary will be available in bin/ensync
```

### 2. Using Go Install
If you have Go installed, you can install EnSync CLI directly using the `go install` command:
```bash
# Install the EnSync CLI tool
go install github.com/rossi1/ensync-cli@latest
```

To install a specific version:
```bash
go install github.com/ensync-cli/cmd/ensync@v1.0.0
```

**Notes**:
1. Requires Go 1.17 or later.
2. Ensure `$GOPATH/bin` (or `$GOBIN`) is in your system’s `PATH`.

---

## Configuration

### Configuration File
Create a config file at `~/.ensync/config.yaml`:
```yaml
base_url: "http://localhost:8080/api/v1/ensync"
api_key: "your-api-key"
debug: false
```

### Environment Variables
```bash
export ENSYNC_API_KEY="your-api-key"
export ENSYNC_BASE_URL="http://localhost:8080/api/v1/ensync"
```

---

## Usage

### Event Management

#### List Events
```bash
ensync event list --page 0 --limit 10 --order DESC --order-by createdAt
ensync event list --order ASC --order-by name
```

#### Create Event
```bash
ensync event create --name "test-event" --payload '{"key":"value","another":"data"}'
```

#### Update Event
```bash
ensync event update --id 123 --name "updated-name"
ensync event update --id 123 --payload '{"key":"new-value"}'
```

### Access Key Management

#### List Access Keys
```bash
ensync access-key list
```

#### Create Access Key
```bash
ensync access-key create --permissions '{"send": ["event1"], "receive": ["event2"]}'
```

#### Manage Permissions
```bash
ensync access-key permissions get --key "MyAccessKey"
ensync access-key permissions set --key "MyAccessKey" --permissions '{"send": ["event1"], "receive": ["event2"]}'
```

### General Options

#### Debug Mode
```bash
ensync --debug event list
```

#### Version Information
```bash
ensync version
ensync version --json
```

---

## Common Flags
- `--page`: Page number for pagination (default: 0)
- `--limit`: Number of items per page (default: 10)
- `--order`: Sort order (ASC/DESC)
- `--order-by`: Field to sort by (name/createdAt)
- `--debug`: Enable debug mode
- `--config`: Specify custom config file location

---

## Error Handling
- Exit code `0`: Success
- Exit code `1`: General error
- Exit code `2`: Configuration error
- Exit code `3`: API error

---

## Development

### Build
```bash
make build
```

### Run Tests
```bash
make test
make test-integration
```

### Lint
```bash
make lint
```

### Release
```bash
make release
```
---