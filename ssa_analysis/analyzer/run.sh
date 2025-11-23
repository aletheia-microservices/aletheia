#!/bin/bash

# Usage:
#   ./run.sh [--eval|--debug]
#
# Examples:
#   ./run.sh --eval
#   ./run.sh --debug
#   ./run.sh

apps=(
    digota
    sockshop3
    postnotification_simple
    dsb_sn2
    dsb_media_nosql
    train_ticket2
    #large_scale_app_A
    #large_scale_app_B
    #large_scale_app_C
    #large_scale_app_D
    #large_scale_app_E
)

mode=""
runs=1

if [[ $# -eq 1 ]]; then
    case "$1" in
        --eval)
            mode="--eval"
            runs=5
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
