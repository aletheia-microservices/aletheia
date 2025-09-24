import matplotlib.pyplot as plt
import numpy as np

apps = ["Digota", "Sockshop", "PostNotification", "SocialNetwork", "MediaMicroservices", "TrainTicket"]
ms_weight = 1
ds_weight = 1
ms_counts = np.array([4, 7, 3, 10, 5, 30])
ds_counts = np.array([5, 6, 2, 10, 4, 21])
complexity = ms_weight * ms_counts + ds_weight * ds_counts

parsing_time   = np.array([10.963, 10.689, 7.861, 17.403, 8.688, 30.954])
schema_time = np.array([0.10, 0.090, 0.007, 0.237, 0.40, 0.371])
detection_time = np.array([0.012, 0.061, 0.005, 0.142, 0.032, 0.376])

#order = np.argsort(complexity)
# use lexsort: primary = complexity, secondary = ms_counts (descending)
order = np.lexsort((ds_counts, complexity))

apps_sorted      = [apps[i] for i in order]
parsing_sorted   = parsing_time[order]
detection_sorted = detection_time[order]
schema_sorted = schema_time[order]

x = np.arange(len(apps_sorted))

plt.figure(figsize=(9,9)) 

parse_color, detect_color = "#80cdc1", "#f4a582"
parse_color, detect_color = "#fdbf6f", "#cab2d6"
parse_color, detect_color = "#a6cee3", "#b2df8a"
parse_color, detect_color = "#9ecae1", "#fdd0a2"
parse_color, detect_color = "#1f77b4", "#ff7f0e"
parse_color, detect_color = "#2ca02c", "#d62728"
parse_color, detect_color = "#8c564b", "#e377c2"
parse_color, detect_color = "#7f7f7f", "#bcbd22"
parse_color, detect_color = "#5db8ea", "#fd9144"
parse_color, schema_color, detect_color = "#5db8ea", "#fd9144", "#9467bd"

# -------------------------------

bars_parse   = plt.bar(x, parsing_sorted, label="Code Parser -> Abstract Graph (Tainted)", color=parse_color)
bars_schema  = plt.bar(x, schema_sorted, bottom=parsing_sorted, label="Schema Builder", color=schema_color)
bars_detect  = plt.bar(x, detection_sorted, bottom=parsing_sorted+schema_sorted, label="Pattern Detector", color=detect_color)

xtick_labels = [f"{app}\n({c})" for app, c in zip(apps_sorted, complexity[order])]
plt.xticks(x, xtick_labels, rotation=20, ha="right")

plt.ylabel("Time (s)")
plt.title("Analysis Time per Application (sorted by complexity)")
plt.legend()

totals = parsing_sorted + detection_sorted
for xi, total in zip(x, totals):
    plt.text(xi, total + max(totals)*0.015, f"{total:.3f}s", ha="center", va="bottom", fontsize=8)

plt.tight_layout()
plt.show()
