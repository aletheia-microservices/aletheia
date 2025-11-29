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
    #sockshop3
    #postnotification_simple
    #dsb_sn2
    #dsb_media_nosql
    #train_ticket2
)

apps_synthetic=(
    #synthetic_app
    #synthetic_appA
    #synthetic_appB
    #synthetic_app1
    #synthetic_app2
    #synthetic_app3
    #synthetic_app4
    #synthetic_app5
    synthetic_app6
)

mode=""
runs=2

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
        go run main.go $mode $app
    done

    echo
done

for app in "${apps_synthetic[@]}"; do
    echo "=== Running $app ($runs times) $mode ==="
    
    for i in $(seq 1 $runs); do
        echo "=== Run $i/$runs"
        go run main.go $mode --synthetic $app
    done

    echo
done
