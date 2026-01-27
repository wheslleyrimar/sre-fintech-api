#!/bin/bash
set -e

# Instala as aplicações na raiz deste projeto (sre).
# Binários gerados em: <raiz-do-sre>/bin/
# fintech-api-failures: clonado de https://github.com/wheslleyrimar/fintech-api-failures

FINTECH_REPO="https://github.com/wheslleyrimar/fintech-api-failures"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
INSTALL_ROOT="$SCRIPT_DIR"
BIN_DIR="$INSTALL_ROOT/bin"
FINTECH_CACHE="${INSTALL_ROOT}/.fintech-api-failures-src"

mkdir -p "$BIN_DIR"

if [ -d "$FINTECH_CACHE/.git" ]; then
  echo "Atualizando fintech-api-failures (repositório remoto)..."
  (cd "$FINTECH_CACHE" && git fetch origin && git reset --hard origin/main)
else
  echo "Clonando fintech-api-failures de $FINTECH_REPO ..."
  rm -rf "$FINTECH_CACHE"
  git clone --depth 1 "$FINTECH_REPO" "$FINTECH_CACHE"
fi

echo "Compilando fintech-api-failures (origem: $FINTECH_CACHE)..."
(cd "$FINTECH_CACHE" && go build -o "$BIN_DIR/fintech-api-failures" ./cmd/api)
echo "Binário instalado: $BIN_DIR/fintech-api-failures (backend na porta 8080)"

echo "Compilando SRE API..."
cd "$INSTALL_ROOT"
go build -o "$BIN_DIR/sre" ./cmd/api
echo "Binário instalado: $BIN_DIR/sre"
echo ""
echo "Todos os binários em: $BIN_DIR"
