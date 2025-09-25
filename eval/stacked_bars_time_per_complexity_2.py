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
total_time_s = parsing_time + schema_time/1000 + detection_time/1000

# sort apps by complexity
order = np.lexsort((ms_counts, complexity))

apps_sorted      = [apps[i] for i in order]
parsing_sorted   = parsing_time[order]
schema_sorted    = schema_time[order]
detection_sorted = detection_time[order]
total_time_s = total_time_s[order]

parsing_time = np.round(parsing_time).astype(int)
total_time_s = np.round(total_time_s).astype(int)

x = np.arange(len(apps_sorted))

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
plt.rcParams['figure.figsize'] = [5, 5]
plt.rcParams['axes.labelsize'] = 'xx-small'
plt.rcParams['legend.fontsize'] = 'xx-small'
plt.rcParams['xtick.labelsize'] = 'xx-small'
plt.rcParams['ytick.labelsize'] = 'xx-small'

fig, axes = plt.subplots(4, 1, sharex=True)

bars0 = axes[0].bar(x, total_time_s, color=COLORS['total'])
axes[0].bar_label(bars0, fmt="%ds", fontsize=6, padding=3)
axes[0].set_title("Total", fontsize=8)
axes[0].set_ylabel("Time (s)")
axes[0].set_ylim(0, 40)

bars1 = axes[1].bar(x, parsing_sorted, color=COLORS['parser'])
axes[1].bar_label(bars1, fmt="%ds", fontsize=6)
axes[1].set_title("Parsing", fontsize=8)
axes[1].set_ylabel("Time (s)")
axes[1].set_ylim(0, 40)

bars2 = axes[2].bar(x, schema_sorted, color=COLORS['schema'])
axes[2].bar_label(bars2, fmt="%dms", fontsize=6)
axes[2].set_title("Schema Building", fontsize=8)
axes[2].set_ylabel("Time (ms)")
axes[2].set_ylim(0, 500)

bars3 = axes[3].bar(x, detection_sorted, color=COLORS['detector'])
axes[3].bar_label(bars3, fmt="%dms", fontsize=6)
axes[3].set_title("Detection", fontsize=8)
axes[3].set_ylabel("Time (ms)")
axes[3].set_ylim(0, 500)

for i, ax in enumerate(axes):
    if i != 0:
      continue

    # secondary y-axis for #ms
    ax2 = ax.twinx()

    # plot #ms
    ax2.plot(
        x, ms_counts[order],
        marker='o', linestyle='-', linewidth=1,
        markersize=3, markerfacecolor='white', markeredgewidth=0.8,
        color='black', label='# microservices', zorder=5
    )
    #ax2.set_ylim(0, 50)
    #ax2.tick_params(axis='y', labelsize=6)

    # plot #ds
    ax2.plot(
        x, ds_counts[order],
        marker='s', linestyle='--', linewidth=1,
        markersize=3, markerfacecolor='white', markeredgewidth=0.8,
        color='black', label='# datastores', zorder=5
    )

    # shared y-axis scaling
    upper_lim = max(ms_counts.max(), ds_counts.max()) * 1.2
    ax2.set_ylim(0, upper_lim)
    ax2.tick_params(axis='y', labelsize=6)
    ax2.set_ylabel('# microservices / \n# datastores', fontsize=7)

    # one combined legend
    ax2.legend(loc='upper left', fontsize=6)
    
    # third y-axis for #ds
    #ax3 = ax.twinx()
    #ax3.spines["right"].set_position(("outward", 35))  # shift outward so it doesn’t overlap ax2
    #ax3.plot(
    #    x, ds_counts[order],
    #    marker='s', linestyle='--', linewidth=1,
    #    markersize=3, markerfacecolor='white', markeredgewidth=0.8,
    #    color='tab:blue', label='# datastores', zorder=5
    #)
    #ax3.set_ylim(0, 25)
    #ax3.tick_params(axis='y', labelsize=6)
#
    #ax2.set_ylabel('# microservices', fontsize=7)
    #ax3.set_ylabel('# datastores', fontsize=7)
    #ax2.legend(loc='upper left', fontsize=6)
    #ax3.legend(loc='upper right', fontsize=6)

    # add labels/legends only on middle subplot to avoid repetition
    #if i == 1:
    #    ax2.set_ylabel('# microservices', fontsize=7)
    #    ax3.set_ylabel('# datastores', fontsize=7)
    #    ax2.legend(loc='upper left', fontsize=6)
    #    ax3.legend(loc='upper right', fontsize=6)
    #else:
    #    ax2.set_ylabel('')
    #    ax3.set_ylabel('')
    #    ax2.set_yticklabels([])
    #    ax3.set_yticklabels([])

# put labels only on the bottom subplot, rotated diagonally
axes[-1].set_xticks(x)
axes[-1].set_xticklabels(xtick_labels, rotation=20, ha="right", rotation_mode="anchor")

# hide x tick labels on upper subplots
for ax in axes[:-1]:
    ax.tick_params(labelbottom=False)

# shared xx
#axes[-1].set_xlabel("Application complexity")

# shared yy (smaller means further away from axis)
#fig.supylabel("Time", fontsize=9, x=0.02)

plt.tight_layout()
plt.subplots_adjust(left=0.12, hspace=0.50)  # smaller hspace = less vertical gap

# smaller means left border
plt.subplots_adjust(left=0.12)
plt.savefig("plot-time-complexity_2.png")
#plt.show()
