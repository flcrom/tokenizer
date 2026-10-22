#!/bin/sh
# @file challenge-cli/testdata/programs/fake_program.sh
# @description Provides deterministic pass, failure, JSON, stdin, file, and timeout behaviors for runner self-tests.
# @depends_on none
# @depended_on_by challenge-cli/internal/runner/runner_test.go

set -eu

mode=${1:-}
case "$mode" in
  echo)
    printf '%s' "$2"
    printf 'diagnostic:%s' "$2" >&2
    ;;
  stdin)
    cat
    ;;
  json)
    printf '{"b":2,"a":1}\n'
    ;;
  fail)
    printf '{"error":"bad"}' >&2
    exit 7
    ;;
  write)
    printf '%s' "$2" > "$3"
    ;;
  timeout)
    exec sleep 5
    ;;
  *)
    printf 'unknown fake mode' >&2
    exit 64
    ;;
esac

