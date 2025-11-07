#!/usr/bin/env bash
# scripts/rebrand.sh
# Used to rebrand the entire repo from dYdX → Vindax
# Usage: ./scripts/rebrand.sh

set -euo pipefail

ROOT="$(git rev-parse --show-toplevel 2>/dev/null || realpath "$(dirname "$0")/..")"
cd "$ROOT"

echo "🚀 Running Go rebrand tool..."
go run ./scripts/rebrand.go

echo "🧹 Cleaning up and verifying..."
go mod tidy

echo "✅ Rebrand complete! Backup stored at: $BACKUP_DIR"
