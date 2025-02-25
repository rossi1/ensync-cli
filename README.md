# EnSync CLI Tool

A command-line interface for interacting with the EnSync config management service.

## Installation

You can install the EnSync CLI Tool in two ways:

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

To install a specific version (e.g., v1.0.1):

```bash
go install github.com/ensync-cli/cmd/ensync@v1.0.1
```

The binary will be installed in your Go binary path (`$GOPATH/bin` or `$GOBIN`).

#### Notes:
1. **Go Version**: `go install` requires Go 1.17 or later.
2. **Path**: After installation, ensure that `$GOPATH/bin` (or `$GOBIN`) is in your system's `PATH` so that you can run the `ensync` command globally.

This method allows you to easily install and update your CLI tool by simply running `go install`, which is especially useful for developers who already have Go set up on their machines.

## Configuration

The CLI can be configured using either a configuration file or environment variables.

### Configuration File
To simplify the configuration, add the environment variables to your shell configuration file:

1. Open your shell configuration file (e.g., `~/.bashrc`, `~/.zshrc`, or `~/.bash_profile`):
   ```bash
   nano ~/.bashrc
   ```

2. Add the following lines:
   ```bash
   export ENSYNC_BASE_URL="http://{url}/api/v1/ensync"
   export ENSYNC_DEBUG=false
   ```

3. Save the file and apply the changes:
   ```bash
   source ~/.bashrc
   ```

### Environment Variables
```bash
export ENSYNC_API_KEY="BjwKUi9EjQtSnR9r9T0MfrrbddIOVCwB"
export ENSYNC_BASE_URL="http://{url}/api/v1/ensync"
```

## Usage

### Event Management

List events:
```bash
# List events with pagination
./bin/ensync --access-key {access-key} event list  --page 0 --limit 10 --order DESC --order-by createdAt

# List events with different ordering
./bin/ensync --access-key {access-key} event list --order ASC --order-by name
```

Create event:
```bash
./bin/ensync event create --access-key {access-key} --name "test-event" --payload '{"key":"value","another":"data"}'
```

Update event:
```bash
# Update event name
./bin/ensync --access-key {access-key} event update  --id 1 --name "updated/name/name"

# Update event payload
./bin/ensync --access-key {your-access-key} event update  --id 1 --payload '{"key":"new-value"}'

# Get event payload by Name
./bin/ensync --access-key {your-access-key} event get --name "updated/name/name"
```

### Access Key Management

List access keys:
```bash
# List all access keys
./bin/ensync --access-key {access-key} access-key list 
```

Create access key:
```bash
# Create access key with permissions
./bin/ensync --access-key {access-key} access-key create --permissions '{"send": ["event1"], "receive": ["event2"]}'
```

Manage permissions:
```bash
# Get current permissions
./bin/ensync --access-key {access-key} access-key permissions get --key {access-key}

# Update permissions
./bin/ensync --access-key {access-key} access-key permissions set --key {access-key} --permissions '{"send": ["event1"], "receive": ["event2"]}'
```

### General Options

Debug mode:
```bash
# Enable debug output
./bin/ensync --debug event list
```

Version information:
```bash
# Display version
ensync version

# Get version in JSON format
ensync version --json
```

## Common Flags

- `--page`: Page number for pagination (default: 0)
- `--limit`: Number of items per page (default: 10)
- `--order`: Sort order (ASC/DESC)
- `--order-by`: Field to sort by (name/createdAt)
- `--debug`: Enable debug mode
- `--config`: Specify custom config file location

## Error Handling

The CLI provides clear error messages and proper exit codes:
- Exit code 0: Success
- Exit code 1: General error
- Exit code 2: Configuration error
- Exit code 3: API error

## Development

Build:
```bash
make build
```

Run tests:
```bash
# Run unit tests
make test

# Run integration tests
make test-integration
```

Lint:
```bash
make lint
```

Release:
```bash
make release
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
