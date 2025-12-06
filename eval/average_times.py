import os
import yaml
from collections import defaultdict
import argparse
from datetime import date
import time

parser = argparse.ArgumentParser()
parser.add_argument("--synthetic", action="store_true", help="enable evaluation mode")
parser.add_argument("--date", type=str, default=str(date.today()), help="date in YYYY-MM-DD format (default: today)")

args = parser.parse_args()

unix_ts = int(time.time())

os.makedirs("results", exist_ok=True)
os.makedirs("results/versions", exist_ok=True)

if args.synthetic:
    INPUT_DIR = f"../ssa_analysis/analyzer/results/times/{args.date}/synthetic"
else:
    INPUT_DIR = f"../ssa_analysis/analyzer/results/times/{args.date}/apps"

if args.synthetic:
    OUTPUT_FILE1 = f"results/averages-synthetic.yaml"
    OUTPUT_FILE2 = f"results/versions/averages-synthetic-{unix_ts}.yaml"
else:
    OUTPUT_FILE1 = f"results/averages-apps.yaml"
    OUTPUT_FILE2 = f"results/versions/averages-apps-{unix_ts}.yaml"

NAME_MAP = {
    "dsb_mediamicroservices":   "MediaMicroservices",
    "dsb_socialnetwork":        "SocialNetwork",
    "postnotification":         "PostNotification",
    "sockshop":                 "SockShop",
    "trainticket":              "TrainTicket",
    "eshopmicroservices":       "EShopMicroservices",
    "digota":                   "Digota",
    "synthetic_app":            "synthetic_app",
    "synthetic_app1":           "App 1",
    "synthetic_app2":           "App 2",
    "synthetic_app3":           "App 3",
    "synthetic_app4":           "App 4",
    "synthetic_app5":           "App 5",
    "synthetic_app6":           "App 6",
}

apps_data = defaultdict(list)
for filename in os.listdir(INPUT_DIR):
    if not filename.endswith(".yaml"):
        continue
    path = os.path.join(INPUT_DIR, filename)

    print(f"[INFO] parsing {path}")

    with open(path, "r") as f:
        data = yaml.safe_load(f)

    if "app" in data:
        key = (data["app"], data["ms_count"], data["ds_count"])
        apps_data[key].append(data)

app_ms_counts = defaultdict(set)
for (app, ms_count, ds_count) in apps_data.keys():
    pretty_app = NAME_MAP.get(app, app)
    app_ms_counts[pretty_app].add(ms_count)

results = {}
for (app, ms_count, ds_count), entries in apps_data.items():
    avg = {}
    n = len(entries)

    for key, value in entries[0].items():
        if key == "app": # skip non-numeric
            continue

        # if the value is numeric, then average it
        if isinstance(value, (int, float)):
            avg[key] = sum(entry[key] for entry in entries) / n
        else:
            # keep the value as-is (e.g., ms_count, ds_count always same)
            avg[key] = value

    # add the app name back
    avg["app"] = app
    avg["ms_count"] = ms_count
    avg["ds_count"] = ds_count
    results[(app, ms_count, ds_count)] = avg

weights = {
    "ms_weight": 1.0,
    "ds_weight": 1.0,
}

ordered_results = []
for (app, ms_count, ds_count), avg in results.items():
    pretty_app = NAME_MAP.get(app, app)

    # if pretty_app has multiple different ms_count values,
    # then label it as "<app>_<ms_count>"
    if len(app_ms_counts[pretty_app]) > 1:
        final_name = f"{pretty_app}_{int(ms_count)}"
    else:
        final_name = pretty_app

    entry = {
        "app":          final_name,
        "iterations":   len(apps_data[(app, ms_count, ds_count)]),
        "ms_count":     int(avg["ms_count"]),
        "ds_count":     int(avg["ds_count"]),
        "rpcs":         int(avg["rpcs"]),
        "total_s":      float(f"{avg['total_s']:.2f}"),
        "parsing_s":    float(f"{avg['parsing_s']:.2f}"),
        "schema_s":     float(f"{avg['schema_s']:.4f}"),
        "detection_s":  float(f"{avg['detection_s']:.4f}"),
    }
    ordered_results.append(entry)

ordered_results = sorted(ordered_results, key=lambda x: x["app"])

with open(OUTPUT_FILE1, "w") as out:
    yaml.dump({"weights": weights}, out, sort_keys=False)
    out.write("\n")

    out.write("apps:\n")
    for entry in ordered_results:
        yaml.dump([entry], out, sort_keys=False)
        out.write("\n")

print(f"[INFO] Saved averaged results to {OUTPUT_FILE1}")

with open(OUTPUT_FILE2, "w") as out:
    yaml.dump({"weights": weights}, out, sort_keys=False)
    out.write("\n")

    out.write("apps:\n")
    for entry in ordered_results:
        yaml.dump([entry], out, sort_keys=False)
        out.write("\n")

print(f"[INFO] Saved averaged results to {OUTPUT_FILE2}")
