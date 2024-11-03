#!/bin/bash

set -euo pipefail

go test -cover -timeout 2s ./...