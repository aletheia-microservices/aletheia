import re
import argparse
import sys
from datetime import date
from pathlib import Path
from collections import defaultdict

PATTERN = re.compile(r"(\d+)\s+maximum resident set size")

parser = argparse.ArgumentParser()
parser.add_argument("--synthetic", action="store_true", help="enable synthetic mode")
parser.add_argument("--date", type=str, default=str(date.today()), help="date in YYYY-MM-DD format (default: today)")
args = parser.parse_args()

INPUT_DIR = Path(f"../ssa_analysis/analyzer/results/metrics/{args.date}")
if args.synthetic:
    INPUT_DIR = INPUT_DIR / "synthetic"

def extract_app_name(filename: str) -> str:
    return filename.split(".", 1)[0]

def extract_peak(path: Path):
    for line in path.read_text(errors="ignore").splitlines():
        m = PATTERN.search(line)
        if m:
            return int(m.group(1))
    return None

def main():
    if not INPUT_DIR.exists():
        print(f"[ERROR] directory not found: {INPUT_DIR}")
        sys.exit(1)

    # group values by app
    groups = defaultdict(list)

    for file in sorted(INPUT_DIR.iterdir()):
        if not file.is_file():
            continue

        peak = extract_peak(file)
        if peak is None:
            continue

        app = extract_app_name(file.name)
        groups[app].append(peak / 1024**2)  # MB

    # compute averages
    averaged = [(app, sum(vals)/len(vals)) for app, vals in groups.items()]

    # print table
    max_app_len = max(len(app) for app, _ in averaged)

    print(f"{'App'.ljust(max_app_len)}   Avg. Peak Memory (MB)")
    print("-" * (max_app_len + 26))

    for app, avg in sorted(averaged, key=lambda x: x[1]):
        print(f"{app.ljust(max_app_len)}   {int(avg)}")

if __name__ == "__main__":
    main()
