# Verify

Runs Aletheia again on each app and checks that it still finds the same number of warnings for each pattern. Use it to confirm that a change doesn't alter the analysis results.

For each app, the script compares the `[NUM_WARNINGS = N]` line of every `output/{app}/analysis/*.txt` with the same file in the baseline, `scripts/verify/expected/{app}/analysis/`. It only compares these counts. It doesn't compare the warnings themselves, or the other files in the baseline (`app.json`, `schema.json`, `ssa/`).

## Running

From the repository root:

```zsh
make verify                                      # every app with a baseline
scripts/verify/verify.sh                         # same as make verify
scripts/verify/verify.sh trainticket sockshop    # only these apps
scripts/verify/verify.sh --config sockshop       # pass config/sockshop.yaml as --detection_config
scripts/verify/verify.sh -h                      # usage
```

The baselines were generated without detection configs. With `--config`, expect lower counts for apps that have a `config/{app}.yaml`.

Apps whose name starts with `synthetic_` are run with `--synthetic`.

## What It Does

1. Builds a new binary into `scripts/verify/logs/{timestamp}/aletheia`, so the run always uses your latest changes.
2. Runs that binary on each app. The results replace `output/{app}/`, and the log goes to `scripts/verify/logs/{timestamp}/{app}.log`. If Aletheia fails, the script puts the previous `output/{app}/` back.
3. Prints the warning count of each pattern, then a summary of all apps.

## Reading the Result

Each app ends with one of:

- `OK`: every count matches the baseline.
- `DIFF`: at least one count differs, shown as `expected N, got M`. The value is `missing` if the file or its `[NUM_WARNINGS = N]` line is missing.
- `ERROR`: Aletheia failed. The last 15 lines of its log are printed.
- `SKIP`: the app has no baseline (only when you name the app explicitly).

The exit code is 0 if no app is `DIFF` or `ERROR`, and 1 otherwise. A failed build also exits with 1.

## Updating the Baseline

If a change is meant to alter the warnings of an app, review the new warnings, then copy them into the baseline:

```zsh
diff -r scripts/verify/expected/sockshop/analysis output/sockshop/analysis
cp output/sockshop/analysis/*.txt scripts/verify/expected/sockshop/analysis/
```

To add a baseline for a new app, analyze it first (`./bin/aletheia {app}`), then:

```zsh
mkdir -p scripts/verify/expected/{app}/analysis
cp output/{app}/analysis/*.txt scripts/verify/expected/{app}/analysis/
```

## Notes

- The script overwrites `output/{app}/` for every app it runs.
- `scripts/verify/logs/` is ignored by git and grows with each run, since every run folder holds a copy of the binary. Delete old runs with `rm -rf scripts/verify/logs/*`.
- The integration tests also compare every app's results with a baseline, `tests/expected/` (see [tests/README.md](../../tests/README.md)).
