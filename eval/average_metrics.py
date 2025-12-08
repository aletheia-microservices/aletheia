import re
import argparse
import sys
import time
from datetime import date
from pathlib import Path
from collections import defaultdict

OUTPUT_DIR_BASE = Path("results")
OUTPUT_DIR_BASE.mkdir(exist_ok=True)
VERSIONS_DIR = OUTPUT_DIR_BASE / "versions"
VERSIONS_DIR.mkdir(exist_ok=True)


PATTERN_MAC = re.compile(r"(\d+)\s+maximum resident set size")
PATTERN_LINUX = re.compile(r"Maximum resident set size \(kbytes\):\s*(\d+)")

parser = argparse.ArgumentParser()
parser.add_argument("--synthetic", action="store_true", help="enable synthetic mode")
parser.add_argument("--date", type=str, default=str(date.today()), help="date in YYYY-MM-DD format (default: today)")
args = parser.parse_args()

INPUT_DIR = Path(f"../ssa-analysis/analyzer/results/metrics/{args.date}")
if args.synthetic:
    INPUT_DIR = INPUT_DIR / "synthetic"
else:
    INPUT_DIR = INPUT_DIR / "apps"

def extract_app_name(filename: str) -> str:
    return filename.split(".", 1)[0]

def extract_peak(path: Path):
    for line in path.read_text(errors="ignore").splitlines():
        m = PATTERN_MAC.search(line)
        if m:
            return int(m.group(1))  # bytes

        m = PATTERN_LINUX.search(line)
        if m:
            return int(m.group(1)) * 1024  # convert kbytes to bytes

    return None

def save(averaged, max_app_len):
    lines = []

    header = f"{'App'.ljust(max_app_len)}   Avg. Peak Memory (MB)"
    sep = "-" * (max_app_len + 26)
    lines.append(header)
    lines.append(sep)

    for app, avg in sorted(averaged, key=lambda x: x[1]):
        lines.append(f"{app.ljust(max_app_len)}   {int(avg)}")

    table_str = "\n".join(lines)

    print(table_str)
    unix_ts = int(time.time())

    filename_base = "metrics-synthetic" if args.synthetic else "metrics-apps"

    with open(OUTPUT_DIR_BASE / f"{filename_base}.txt", "w") as f:
        f.write(table_str + "\n")

    with open(VERSIONS_DIR / f"{filename_base}-{unix_ts}.txt", "w") as f:
        f.write(table_str + "\n")

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

    averaged = [(app, sum(vals)/len(vals)) for app, vals in groups.items()]
    max_app_len = max(len(app) for app, _ in averaged)
    save(averaged, max_app_len)

if __name__ == "__main__":
    main()
