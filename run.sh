#!/bin/sh

set -e 

(
  cd "$(dirname "$0")" 
  go build -o /tmp/interpreter-target ./cmd
)

exec /tmp/interpreter-target "$@"
