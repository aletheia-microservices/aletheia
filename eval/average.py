import os
import yaml
from collections import defaultdict

INPUT_DIR = "../ssa_analysis/analyzer/analysis_times"
OUTPUT_FILE = "results/averages.yaml"

NAME_MAP = {
    "dsb_media_nosql":          "mediamicroservices",
    "dsb_sn2":                  "socialnetwork",
    "large_scale_app":          "largescaleapp",
    "postnotification_simple":  "postnotification",
    "sockshop3":                "sockshop",
    "train_ticket2":            "trainticket",
}

apps_data = defaultdict(list)
for filename in os.listdir(INPUT_DIR):
    if not filename.endswith(".yaml"):
        continue
    path = os.path.join(INPUT_DIR, filename)

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
        "total_s":      int(avg["total_s"]),
        "init_s":       int(avg["init_s"]),
        "parsing_s":    int(avg["parsing_s"]),
        "schema_s":     float(f"{avg['schema_s']:.3f}"),
        "detection_s":  float(f"{avg['detection_s']:.3f}"),
    }
    ordered_results.append(entry)

ordered_results = sorted(ordered_results, key=lambda x: x["app"])

with open(OUTPUT_FILE, "w") as out:
    yaml.dump({"weights": weights}, out, sort_keys=False)
    out.write("\n")

    out.write("apps:\n")
    for entry in ordered_results:
        yaml.dump([entry], out, sort_keys=False)
        out.write("\n")

print(f"[INFO] Saved averaged results to {OUTPUT_FILE}")
