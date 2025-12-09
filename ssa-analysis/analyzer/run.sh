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
)

mode=""
runs=1

DATE=$(date +%F)
METRICS_DIR="results/metrics/$DATE"
mkdir -p "$METRICS_DIR"
mkdir -p "$METRICS_DIR/apps"
mkdir -p "$METRICS_DIR/synthetic"

if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS uses BSD time with -l
    TIME_CMD="/usr/bin/time -l"
else
    # Linux uses GNU time with -v
    TIME_CMD="time -v"
fi

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

# https://man.freebsd.org/cgi/man.cgi?manpath=macOS+13.6.5&query=getrusage&sektion=2
# https://stackoverflow.com/questions/41205680/how-to-get-the-memory-usage-of-a-os-x-macos-process
# https://unix.stackexchange.com/questions/77370/how-to-measure-on-linux-the-peak-memory-of-an-application-after-has-ended
# https://stackoverflow.com/questions/774556/peak-memory-usage-of-a-linux-unix-process
# macOS: /usr/bin/time (BSD) with -l -> resource usage
# time writes to stderr, so redirect 2> to capture metrics

for app in "${apps[@]}"; do
    echo "=== Running $app ($runs times) $mode ==="

    for i in $(seq 1 $runs); do
        echo "=== Run $i/$runs"
        timestamp=$(date +%s)
        output_file="$METRICS_DIR/apps/${app}.${timestamp}.txt"
        timestamp=$(date +%s)
        output_file="$METRICS_DIR/apps/${app}.${timestamp}.txt"
        if [[ "$mode" == "--eval" ]]; then
            $TIME_CMD go run main.go $mode "$app" 2> "$output_file"
        else
            go run main.go $mode "$app"
        fi
    done

    echo
done

for app in "${apps_synthetic[@]}"; do
    echo "=== Running $app ($runs times) $mode ==="

    # choose build tags based on synthetic app
    go_tags=""
    case "$app" in
        synthetic_app1|synthetic_app2|synthetic_app3)
            go_tags="-tags synthetic_small"
            ;;
        synthetic_app4)
            go_tags="-tags synthetic_medium"
            ;;
        synthetic_app5)
            go_tags="-tags synthetic_large"
            ;;
    esac

    for i in $(seq 1 $runs); do
        echo "=== Run $i/$runs"
        timestamp=$(date +%s)
        output_file="$METRICS_DIR/synthetic/${app}.${timestamp}.txt"

        if [[ "$mode" == "--eval" ]]; then
            $TIME_CMD go run $go_tags main.go $mode --synthetic "$app" 2> "$output_file"
        else
            go run $go_tags main.go $mode --synthetic "$app"
        fi
    done

    echo
done
