#!/usr/bin/env bash
set -euo pipefail

target="linux"
skip_docker=0
skip_build=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    linux|windows)
      target="$1"
      shift
      ;;
    --target)
      target="${2:-}"
      shift 2
      ;;
    --skip-docker)
      skip_docker=1
      shift
      ;;
    --skip-build)
      skip_build=1
      shift
      ;;
    -h|--help)
      echo "Usage: ./run.sh [linux|windows] [--target linux|windows] [--skip-docker] [--skip-build]"
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

if [[ "$target" != "linux" && "$target" != "windows" ]]; then
  echo "Target must be linux or windows." >&2
  exit 1
fi

run_step() {
  local name="$1"
  shift

  echo
  echo "==> $name"
  "$@"
}

if [[ "$skip_docker" -eq 0 ]]; then
  run_step "Starting Grafana with Docker Compose" docker compose up -d --build
fi

if [[ "$skip_build" -eq 0 ]]; then
  run_step "Building backend for $target" mage -v "build:$target"
fi

run_step "Starting frontend watcher" npm run dev
