#!/usr/bin/env bash
# Re-runs Aletheia and compares the [NUM_WARNINGS = N] of each analysis/*.txt in output/
# against the baseline in scripts/verify/expected/ (logs are saved in scripts/verify/logs/).
#
# usage: scripts/verify/verify.sh [--config] [app ...]
#   app ...   apps to check (default: every app with a baseline in scripts/verify/expected/)
#   --config  pass config/{app}.yaml as --detection_config when it exists (the baseline was generated without it)
#
# Exit code is 0 if every app matches its baseline, 1 otherwise.

set -o pipefail

BASELINE=scripts/verify/expected
OUTPUT=output
USE_CONFIG=false

apps=()
for arg in "$@"; do
    case "$arg" in
        --config) USE_CONFIG=true ;;
        -h|--help) sed -n '2,/^$/s/^# \{0,1\}//p' "$0"; exit 0 ;;
        *) apps+=("$arg") ;;
    esac
done

cd "$(dirname "$0")/../.."

if [ ${#apps[@]} -eq 0 ]; then
    for d in "$BASELINE"/*/; do
        [ -d "$d" ] && apps+=("$(basename "$d")")
    done
fi

if [ ${#apps[@]} -eq 0 ]; then
    echo "no apps found in $BASELINE/"
    exit 1
fi

GREEN="" YELLOW="" RED="" BLUE="" RESET=""
if [ -t 1 ]; then
    GREEN=$'\e[32m' YELLOW=$'\e[33m' RED=$'\e[31m' BLUE=$'\e[34m' RESET=$'\e[0m'
fi

# prints its arguments in the given color
say() {
    local color=$1; shift
    echo "${color}$*${RESET}"
}

# prints the N in [NUM_WARNINGS = N] of the given file, or "missing" if there is none
num_warnings() {
    local n
    n=$(sed -n 's/^\[NUM_WARNINGS = \([0-9]*\)\]$/\1/p' "$1" 2>/dev/null | head -n 1)
    echo "${n:-missing}"
}

# compares the warning counts of every analysis file in the baseline or output of the given app
compare_warnings() {
    local app=$1 name expected actual rc=0
    for name in $(for f in "$BASELINE/$app"/analysis/*.txt "$OUTPUT/$app"/analysis/*.txt; do
                      [ -f "$f" ] && basename "$f" .txt
                  done | sort -u); do
        expected=$(num_warnings "$BASELINE/$app/analysis/$name.txt")
        actual=$(num_warnings "$OUTPUT/$app/analysis/$name.txt")
        if [ "$expected" = "$actual" ]; then
            printf '    %-28s %s\n' "$name" "$actual"
        else
            say "$YELLOW" "$(printf '    %-28s expected %s, got %s  <--' "$name" "$expected" "$actual")"
            rc=1
        fi
    done
    return $rc
}

LOGDIR="scripts/verify/logs/$(date +%Y%m%d-%H%M%S)"
mkdir -p "$LOGDIR"
BIN="$LOGDIR/aletheia"

echo "building aletheia..."
if ! go build -o "$BIN" ./cmd/aletheia; then
    say "$RED" "build failed"
    exit 1
fi

passed=() failed=() errored=() skipped=()

for app in "${apps[@]}"; do
    echo
    echo "=== $app ==="

    if [ ! -d "$BASELINE/$app" ]; then
        say "$BLUE" "SKIP: no baseline in $BASELINE/$app"
        skipped+=("$app")
        continue
    fi

    args=()
    $USE_CONFIG && [ -f "config/$app.yaml" ] && args+=(--detection_config "config/$app.yaml")
    [[ "$app" == synthetic_* ]] && args+=(--synthetic)

    # start from a clean directory so stale files from previous runs cannot hide differences,
    # keeping the previous output aside so it can be restored if the run fails
    prev="$LOGDIR/prev-$app"
    [ -d "$OUTPUT/$app" ] && mv "$OUTPUT/$app" "$prev"

    log="$LOGDIR/$app.log"
    echo "running: aletheia ${args[*]} $app  (log: $log)"
    if ! "$BIN" "${args[@]}" "$app" >"$log" 2>&1; then
        say "$RED" "ERROR: aletheia failed, last lines of log:"
        tail -n 15 "$log" | sed 's/^/    /'
        if [ -d "$prev" ]; then
            rm -rf "${OUTPUT:?}/$app"
            mv "$prev" "$OUTPUT/$app"
            echo "restored previous $OUTPUT/$app"
        fi
        errored+=("$app")
        continue
    fi
    rm -rf "$prev"

    if compare_warnings "$app"; then
        say "$GREEN" "OK: warning counts match baseline"
        passed+=("$app")
    else
        say "$YELLOW" "DIFF: warning counts differ from baseline"
        failed+=("$app")
    fi
done

echo
echo "=== summary ==="
say "$GREEN"  "passed:  ${#passed[@]}  ${passed[*]:-}"
say "$YELLOW" "differ:  ${#failed[@]}  ${failed[*]:-}"
say "$RED"    "errored: ${#errored[@]}  ${errored[*]:-}"
say "$BLUE"   "skipped: ${#skipped[@]}  ${skipped[*]:-}"
echo "logs: $LOGDIR"

[ ${#failed[@]} -eq 0 ] && [ ${#errored[@]} -eq 0 ]
