import matplotlib.pyplot as plt
import seaborn as sns
import numpy as np
import yaml
import argparse
import time
import os

parser = argparse.ArgumentParser()
parser.add_argument("--synthetic", action="store_true", help="enable synthetic app")
args = parser.parse_args()

unix_ts = int(time.time())

os.makedirs("plots", exist_ok=True)
os.makedirs("plots/versions", exist_ok=True)

if args.synthetic:
  INPUT_FILE = "results/averages-synthetic.yaml"
else:
  INPUT_FILE = "results/averages-apps.yaml"

if args.synthetic:
  OUTPUT_FILE1 = f"plots/plot-synthetic.png"
  OUTPUT_FILE2 = f"plots/versions/plot-synthetic-{unix_ts}.png"
else:
  OUTPUT_FILE1 = f"plots/plot-apps.png"
  OUTPUT_FILE2 = f"plots/versions/plot-apps-{unix_ts}.png"

# load yaml data
with open(INPUT_FILE, "r") as f:
  data = yaml.safe_load(f)

# weights
ms_weight = data["weights"]["ms_weight"]
ds_weight = data["weights"]["ds_weight"]

# actual results
apps = [app["app"] for app in data["apps"]]
rpcs = np.array([app["rpcs"] for app in data["apps"]])
ms_counts = np.array([app["ms_count"] for app in data["apps"]])
ds_counts = np.array([app["ds_count"] for app in data["apps"]])
parsing_s = np.array([app["parsing_s"] for app in data["apps"]])
schema_s = np.array([app["schema_s"] for app in data["apps"]])
detection_s = np.array([app["detection_s"] for app in data["apps"]])

# x and y values
complexity = ms_weight * ms_counts + ds_weight * ds_counts
total_s = parsing_s + schema_s + detection_s

print("Apps:", apps)
print("Complexity:", complexity)
print("Total time (s):", total_s)

if args.synthetic:
  order = np.arange(len(ms_counts))
else:
  order = np.lexsort((ms_counts, complexity))

apps_sorted      = [apps[i] for i in order]
parsing_sorted   = parsing_s[order]
schema_sorted    = schema_s[order]
detection_sorted = detection_s[order]
total_s = total_s[order]

spacing = 0.8  # smaller => bars closer horizontally
x = np.arange(len(apps_sorted)) * spacing
bar_width = 0.6

BASE_COLOR_PALETTE = sns.color_palette('deep', 12)
COLORS = {
  'parser': BASE_COLOR_PALETTE[0],
  'schema': BASE_COLOR_PALETTE[1],
  'detector': BASE_COLOR_PALETTE[3],
  'total': BASE_COLOR_PALETTE[2],
}

xtick_labels = [f"{app}" for app, c in zip(apps_sorted, complexity[order])]

sns.set_theme(style='ticks')
plt.rcParams['figure.dpi'] = 600
plt.rcParams['figure.figsize'] = [4, 5]
plt.rcParams['axes.labelsize'] = 'xx-small'
plt.rcParams['legend.fontsize'] = 'xx-small'
plt.rcParams['xtick.labelsize'] = 'xx-small'
plt.rcParams['ytick.labelsize'] = 'xx-small'

fig, axes = plt.subplots(4, 1, sharex=True)

bars0 = axes[0].bar(x, total_s, color=COLORS['total'], width=bar_width)
axes[0].margins(y=0.3) # 20% vertical padding
axes[0].bar_label(bars0, fmt="%.2fs", fontsize=5.5, padding=3)
axes[0].set_title("Total", fontsize=8)

bars1 = axes[1].bar(x, parsing_sorted, color=COLORS['parser'], width=bar_width)
axes[1].margins(y=0.2) # 20% vertical padding
axes[1].bar_label(bars1, fmt="%.2fs", fontsize=5.5)
axes[1].set_title("Parsing", fontsize=8)

bars2 = axes[2].bar(x, schema_sorted, color=COLORS['schema'], width=bar_width)
axes[2].margins(y=0.2) # 20% vertical padding
axes[2].bar_label(bars2, fmt="%.4fs", fontsize=5.5)
axes[2].set_title("Schema Building", fontsize=8)

bars3 = axes[3].bar(x, detection_sorted, color=COLORS['detector'], width=bar_width)
axes[3].margins(y=0.2) # 20% vertical padding
axes[3].bar_label(bars3, fmt="%.4fs", fontsize=5.5)
axes[3].set_title("Detection", fontsize=8)

for i, ax in enumerate(axes):
    if i != 0:
      continue
    # secondary y-axis for #ms
    if args.synthetic:
      ax2 = ax.twinx()
      ax2.set_yscale("log")
    else:
      ax2 = ax.twinx()
    # plot #ms
    ax2.plot(
        x, ms_counts[order],
        marker='o', linestyle='-', linewidth=1,
        markersize=3, markerfacecolor='white', markeredgewidth=0.8,
        color='black', label='# microservices', zorder=5
    )
    # plot #ds
    ax2.plot(
        x, ds_counts[order],
        marker='s', linestyle='--', linewidth=1,
        markersize=3, markerfacecolor='white', markeredgewidth=0.8,
        color='black', label='# datastores', zorder=5
    )
    if args.synthetic:
      # plot #rpcs
      ax2.plot(
          x, rpcs[order],
          marker='^', linestyle='--', linewidth=1,
          markersize=3, markerfacecolor='white', markeredgewidth=0.8,
          color='black', label='# rpcs', zorder=5
    )
    # shared y-axis scaling
    if args.synthetic:
      upper_lim = max(ms_counts.max(), ds_counts.max(), rpcs.max()) * 2.2
    else:
      upper_lim = max(ms_counts.max(), ds_counts.max(), rpcs.max()) * 1.2
    ax2.set_ylim(0, upper_lim)
    ax2.tick_params(axis='y', labelsize=6)
    ax2.set_ylabel('# microservices / \n# datastores / # rpcs', fontsize=7)
    # one combined legend
    ax2.legend(loc='upper left', fontsize=6)

# shared y-axis label
fig.text(0.02, 0.5, "Time (s)", va='center', rotation='vertical', fontsize=8)

# put labels only on the bottom subplot, rotated diagonally
axes[-1].set_xticks(x)
axes[-1].set_xticklabels(xtick_labels, rotation=35, ha="right", rotation_mode="anchor")

# hide x tick labels on upper subplots
for ax in axes[:-1]:
    ax.tick_params(labelbottom=False)

if args.synthetic:
  axes[0].set_yscale("log")
  axes[1].set_yscale("log")
  axes[2].set_yscale("log")
  axes[3].set_yscale("log")

plt.tight_layout()
# smaller hspace => less vertical gap
plt.subplots_adjust(left=0.12, hspace=0.40)
# smaller left => left border
plt.subplots_adjust(left=0.12)
plt.savefig(OUTPUT_FILE1)
print(f"[INFO] saved plot to {OUTPUT_FILE1}")
plt.savefig(OUTPUT_FILE2)
print(f"[INFO] saved plot to {OUTPUT_FILE2}")
