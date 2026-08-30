#!/usr/bin/env bash
set -euo pipefail

[ "$(git rev-parse --abbrev-ref HEAD)" = "main" ] || exit 0

printf '\033[33mCommit to main is blocked. Create a branch first.\033[0m\n' >&2
exit 1
