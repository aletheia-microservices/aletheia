#!/bin/bash

# Usage:
#   ./run.sh [--eval|--debug]
#
# Examples:
#   ./run.sh --eval
#   ./run.sh --debug
#   ./run.sh

set -e

apps=(
    #digota
    #sockshop
    #eshopmicroservices
    #postnotification
    #dsb_socialnetwork
    #dsb_mediamicroservices
    #trainticket
)

apps_synthetic=(
    synthetic_app1
    synthetic_app2
    synthetic_app3
    synthetic_app4
    synthetic_app5
    synthetic_app6
)

mode=""
runs=1

METRICS_DIR="metrics"
mkdir -p "$METRICS_DIR"

if [[ $# -eq 1 ]]; then
    case "$1" in
        --eval)
            mode="--eval"
            ;;
        --debug)
            mode=""
            runs=1
            ;;
        *)
            echo "[ERROR] unknown argument: $1"
            exit 1
            ;;
    esac
fi

for app in "${apps[@]}"; do
    echo "=== Running $app ($runs times) $mode ==="

    for i in $(seq 1 $runs); do
        echo "=== Run $i/$runs"

        # macOS: /usr/bin/time (BSD) with -l -> resource usage
        # time writes to stderr, so redirect 2> to capture metrics
        /usr/bin/time -l \
            go run main.go $mode "$app" \
            2> "$METRICS_DIR/${app}_run${i}.txt"
    done

    echo
done

for app in "${apps_synthetic[@]}"; do
    echo "=== Running $app ($runs times) $mode ==="

    for i in $(seq 1 $runs); do
        echo "=== Run $i/$runs"

        /usr/bin/time -l \
            go run main.go $mode --synthetic "$app" \
            2> "$METRICS_DIR/${app}_synthetic_run${i}.txt"
    done

    echo
done
