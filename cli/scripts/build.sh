#!/usr/bin/env bash
# Cross-compile the ovl CLI for linux and macOS (amd64 + arm64).
# Output binaries are written to build/ at the repo root.
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

OUT_DIR="../build"
mkdir -p "$OUT_DIR"

targets=(
  "linux amd64"
  "linux arm64"
  "darwin amd64"
  "darwin arm64"
)

for target in "${targets[@]}"; do
  read -r os arch <<<"$target"
  name="ovl-${os}-${arch}"
  echo "Building ${name}..."
  GOOS="$os" GOARCH="$arch" mise exec -- go build -o "${OUT_DIR}/${name}" .
done

echo "Done. Binaries written to ${OUT_DIR}/"
