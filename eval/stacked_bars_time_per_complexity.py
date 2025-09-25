import matplotlib.pyplot as plt
import seaborn as sns
import numpy as np

apps = ["Digota", "Sockshop", "PostNotification", "SocialNetwork", "MediaMicroservices", "TrainTicket"]
ms_weight = 1
ds_weight = 1
ms_counts = np.array([4, 7, 3, 10, 5, 30])
ds_counts = np.array([5, 6, 2, 10, 4, 21])
complexity = ms_weight * ms_counts + ds_weight * ds_counts

#parsing_time   = np.array([10.963, 10.689, 7.861, 17.403, 8.688, 30.954])
#schema_time = np.array([0.10, 0.090, 0.007, 0.237, 0.40, 0.371])
#detection_time = np.array([0.012, 0.061, 0.005, 0.142, 0.032, 0.376])

#order = np.argsort(complexity)
# use lexsort: primary = complexity, secondary = ms_counts (descending)

parsing_time   = np.array([10.96, 10.69, 7.86, 17.40, 8.69, 30.95]) #s
schema_time    = np.array([100, 90, 7, 237, 400, 371]) #ms
detection_time = np.array([12, 61, 5, 142, 32, 376]) #ms

parsing_time = np.round(parsing_time).astype(int)

# sort apps by complexity
order = np.lexsort((ms_counts, complexity))

apps_sorted      = [apps[i] for i in order]
parsing_sorted   = parsing_time[order]
schema_sorted    = schema_time[order]
detection_sorted = detection_time[order]

x = np.arange(len(apps_sorted))

BASE_COLOR_PALETTE = sns.color_palette('deep', 12)
COLORS = {
  'parser': BASE_COLOR_PALETTE[0],
  'schema': BASE_COLOR_PALETTE[1],
  'detector': BASE_COLOR_PALETTE[2],
}

xtick_labels = [f"{app}" for app, c in zip(apps_sorted, complexity[order])]

sns.set_theme(style='ticks')
plt.rcParams['figure.dpi'] = 600
plt.rcParams['figure.figsize'] = [5, 5]
plt.rcParams['axes.labelsize'] = 'xx-small'
plt.rcParams['legend.fontsize'] = 'xx-small'
plt.rcParams['xtick.labelsize'] = 'xx-small'
plt.rcParams['ytick.labelsize'] = 'xx-small'

fig, axes = plt.subplots(3, 1, sharex=True)

bars0 = axes[0].bar(x, parsing_sorted, color=COLORS['parser'])
axes[0].bar_label(bars0, fmt="%ds", fontsize=6)
axes[0].set_title("Parsing", fontsize=8)
axes[0].set_ylabel("Time (s)")
axes[0].set_ylim(0, 40)

bars1 = axes[1].bar(x, schema_sorted, color=COLORS['schema'])
axes[1].bar_label(bars1, fmt="%dms", fontsize=6)
axes[1].set_title("Schema Building", fontsize=8)
axes[1].set_ylabel("Time (ms)")
axes[1].set_ylim(0, 500)

bars2 = axes[2].bar(x, detection_sorted, color=COLORS['detector'])
axes[2].bar_label(bars2, fmt="%dms", fontsize=6)
axes[2].set_title("Detection", fontsize=8)
axes[2].set_ylabel("Time (ms)")
axes[2].set_ylim(0, 500)

plt.xticks(x, xtick_labels, rotation=20, ha="right")

# shared xx
#axes[-1].set_xlabel("Application complexity")

# shared yy (smaller means further away from axis)
#fig.supylabel("Time", fontsize=9, x=0.02)

plt.tight_layout()

# smaller means left border
plt.subplots_adjust(left=0.12)
plt.savefig("plot-time-complexity.png")
#plt.show()
