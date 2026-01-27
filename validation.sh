#!/bin/bash
set -e

ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

usage() {
  echo "usage: $0 <app_name> <group_number:1|2|3|4|5|6|7|8> <env:local> <case:1|2|3|4>" >&2
  echo "  env must be 'local'. Backend runs on 8080, SRE on 8081." >&2
  exit 1
}

BACKEND_PORT=8080
SRE_PORT=8081
BIN_DIR="$ROOT/bin"
FINTECH_BIN="$BIN_DIR/fintech-api-failures"
SRE_BIN="$BIN_DIR/sre"

setup() {
  if [ -n "$FINTECH_BINARY" ]; then
    "$FINTECH_BINARY" > backend.log 2>&1 &
  elif [ -x "$FINTECH_BIN" ]; then
    "$FINTECH_BIN" > backend.log 2>&1 &
  else
    echo "Binário do backend não encontrado: $FINTECH_BIN"
    echo "Rode ./install.sh para gerar os binários em $BIN_DIR"
    exit 1
  fi
  while ! nc -z localhost $BACKEND_PORT 2>/dev/null; do
    echo "waiting for backend to start..."
    sleep 0.25
  done
  if [ ! -x "$SRE_BIN" ]; then
    echo "Binário SRE não encontrado: $SRE_BIN"
    echo "Rode ./install.sh para gerar os binários em $BIN_DIR"
    exit 1
  fi
  BACKEND_URL="http://localhost:$BACKEND_PORT" SRE_BASE_URL="http://localhost:$SRE_PORT" PORT=$SRE_PORT "$SRE_BIN" > app.log 2>&1 &
  while ! nc -z localhost $SRE_PORT 2>/dev/null; do
    echo "waiting for SRE application to start..."
    sleep 0.25
  done
}

tearDown() {
  for port in $BACKEND_PORT $SRE_PORT; do
    while nc -z localhost $port 2>/dev/null; do
      lsof -i :$port 2>/dev/null | tail -n1 | awk '{print $2}' | xargs kill 2>/dev/null || true
      sleep 0.1
    done
  done
}

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ] || [ -z "$4" ]; then
  usage
fi

app=$1
group=$2
env=$3
case_number=$4

if [ "$env" != "local" ]; then
  echo "Only env=local is supported."
  usage
fi

echo "running validation for" "$app" "$group" "$env" "$case_number"

export APP_NAME=$app
export GROUP_NUMBER=$group
export ENV=$env
export SRE_URL="http://localhost:$SRE_PORT"

tearDown
setup
trap tearDown EXIT
k6 run "validations/case_${case_number}.js"
tearDown
