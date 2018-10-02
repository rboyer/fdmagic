#!/bin/bash

set -euo pipefail

cd "$(dirname "$0")"

(
    cd child
    go build
)

exec go run parent.go

