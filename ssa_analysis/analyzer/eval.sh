#!/bin/bash

apps=(
    digota
    sockshop3
    postnotification_simple
    dsb_sn2
    dsb_media_nosql
    train_ticket2
    large_scale_app
)

runs=5

for app in "${apps[@]}"; do
    echo "=== Running $app ($runs times) ==="
    
    for i in $(seq 1 $runs); do
        echo "=== Run $i/$runs"
        go run main.go --eval "$app"
    done

    echo
done
