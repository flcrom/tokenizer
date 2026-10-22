#!/bin/sh
# @file your_program.sh
# @description Compiles and executes the learner tokenizer through the stable black-box challenge interface.
# @depends_on go.mod, app/main.go
# @depended_on_by challenge-cli/internal/runner/runner.go, challenge-cli/README.md, docs/harness-architecture.md

set -eu

project_root=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
invocation_dir=$(pwd)
binary_dir="$project_root/.local/bin"
binary_path="$binary_dir/tokenizer"
go_cache="$project_root/.local/go-build"
go_module_cache="$project_root/.local/go-mod"

mkdir -p "$binary_dir" "$go_cache" "$go_module_cache"
(
  cd "$project_root"
  GOCACHE="${GOCACHE:-$go_cache}" GOMODCACHE="${GOMODCACHE:-$go_module_cache}" go build -o "$binary_path" ./app
)

cd "$invocation_dir"
exec "$binary_path" "$@"

