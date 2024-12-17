
#!/bin/sh

set -ex  # Print commands and exit on errors

# Change to the script's directory
cd "$(dirname "$0")"

# Build the Go binary and ensure build errors are shown
go build -o /tmp/interpreter-target ./cmd

# Verify that the build succeeded
if [ ! -f /tmp/interpreter-target ]; then
  echo "Build failed or /tmp/interpreter-target does not exist"
  exit 1
fi

# Execute the binary with the provided arguments
exec /tmp/interpreter-target "$@"

