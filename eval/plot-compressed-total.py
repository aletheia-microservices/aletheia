import matplotlib.pyplot as plt
import seaborn as sns
import numpy as np
import yaml
import time
import os

unix_ts = int(time.time())

os.makedirs("plots", exist_ok=True)
os.makedirs("plots/versions", exist_ok=True)

INPUT_APPS      = "results/averages-apps.yaml"
INPUT_SYNTHETIC = "results/averages-synthetic.yaml"

OUTPUT_FILE1 = "plots/plot-apps-vs-synthetic.png"
OUTPUT_FILE2 = f"plots/versions/plot-apps-vs-synthetic-{unix_ts}.png"

# ------------ load yaml data ------------
with open(INPUT_APPS, "r") as f:
    data_apps = yaml.safe_load(f)

with open(INPUT_SYNTHETIC, "r") as f:
    data_syn = yaml.safe_load(f)

# assume same weights in both (or only use from apps)
ms_weight = data_apps["weights"]["ms_weight"]
ds_weight = data_apps["weights"]["ds_weight"]

# =====================================================
# ------------ REAL APPS DATA (subplot 0) -------------
# =====================================================
apps_real        = [app["app"]        for app in data_apps["apps"]]
rpcs_real        = np.array([app["rpcs"]        for app in data_apps["apps"]])
ms_counts_real   = np.array([app["ms_count"]    for app in data_apps["apps"]])
ds_counts_real   = np.array([app["ds_count"]    for app in data_apps["apps"]])
parsing_real     = np.array([app["parsing_s"]   for app in data_apps["apps"]])
schema_real      = np.array([app["schema_s"]    for app in data_apps["apps"]])
detection_real   = np.array([app["detection_s"] for app in data_apps["apps"]])

complexity_real  = ms_weight * ms_counts_real + ds_weight * ds_counts_real
total_real       = parsing_real + schema_real + detection_real

print("Real Apps:", apps_real)
print("Real Complexity:", complexity_real)
print("Real Total time (s):", total_real)

# sort real apps (same as before)
order_real = np.lexsort((ms_counts_real, total_real))

apps_real_sorted       = [apps_real[i] for i in order_real]
parsing_real_sorted    = parsing_real[order_real]
schema_real_sorted     = schema_real[order_real]
detection_real_sorted  = detection_real[order_real]
total_real_sorted      = total_real[order_real]
rpcs_real_sorted       = rpcs_real[order_real]
ms_counts_real_sorted  = ms_counts_real[order_real]
ds_counts_real_sorted  = ds_counts_real[order_real]

# =====================================================
# ---------- SYNTHETIC APPS DATA (subplot 1) ----------
# =====================================================
apps_syn        = [app["app"]        for app in data_syn["apps"]]
rpcs_syn        = np.array([app["rpcs"]        for app in data_syn["apps"]])
ms_counts_syn   = np.array([app["ms_count"]    for app in data_syn["apps"]])
ds_counts_syn   = np.array([app["ds_count"]    for app in data_syn["apps"]])
parsing_syn     = np.array([app["parsing_s"]   for app in data_syn["apps"]])
schema_syn      = np.array([app["schema_s"]    for app in data_syn["apps"]])
detection_syn   = np.array([app["detection_s"] for app in data_syn["apps"]])

complexity_syn  = ms_weight * ms_counts_syn + ds_weight * ds_counts_syn
total_syn       = parsing_syn + schema_syn + detection_syn

print("Synthetic Apps:", apps_syn)
print("Synthetic Complexity:", complexity_syn)
print("Synthetic Total time (s):", total_syn)

# if you want synthetic also sorted by total time:
order_syn = np.lexsort((ms_counts_syn, total_syn))
# or keep in original order: order_syn = np.arange(len(ms_counts_syn))

apps_syn_sorted       = [apps_syn[i] for i in order_syn]
parsing_syn_sorted    = parsing_syn[order_syn]
schema_syn_sorted     = schema_syn[order_syn]
detection_syn_sorted  = detection_syn[order_syn]
total_syn_sorted      = total_syn[order_syn]
rpcs_syn_sorted       = rpcs_syn[order_syn]
ms_counts_syn_sorted  = ms_counts_syn[order_syn]
ds_counts_syn_sorted  = ds_counts_syn[order_syn]

# =====================================================
# ----------------- plotting config -------------------
# =====================================================
spacing   = 0.8  # smaller => bars closer horizontally
bar_width = 0.6

# separate x axes per subplot
x_real = np.arange(len(apps_real_sorted)) * spacing
x_syn  = np.arange(len(apps_syn_sorted))  * spacing

BASE_COLOR_PALETTE = sns.color_palette('deep', 12)
COLORS = {
    'parser':       BASE_COLOR_PALETTE[0],
    'schema':       BASE_COLOR_PALETTE[1],
    'detector':     BASE_COLOR_PALETTE[3],
    'total':        BASE_COLOR_PALETTE[2],
    'total_real':   BASE_COLOR_PALETTE[0],
    'total_syn':    BASE_COLOR_PALETTE[1],
}

xtick_labels_real = apps_real_sorted
xtick_labels_syn  = apps_syn_sorted

sns.set_theme(style='ticks')
plt.rcParams['figure.dpi']        = 600
plt.rcParams['figure.figsize']    = [4, 4]
plt.rcParams['axes.labelsize']    = 'xx-small'
plt.rcParams['legend.fontsize']   = 'xx-small'
plt.rcParams['xtick.labelsize']   = 'xx-small'
plt.rcParams['ytick.labelsize']   = 'xx-small'

fig, axes = plt.subplots(2, 1)

# =================== subplot 0: REAL APPS ===================
bars0 = axes[0].bar(x_real, total_real_sorted, color=COLORS['total_real'], width=bar_width)
axes[0].margins(y=0.3)  # 30% vertical padding
axes[0].bar_label(bars0, fmt="%.2fs", fontsize=5.5, padding=3)
axes[0].set_title("Total (realistic applications)", fontsize=6, pad=2)

ax0_2 = axes[0].twinx()
# overlay #ms / #ds / #rpcs for real apps (linear scale)
ax0_2.plot(
    x_real, ms_counts_real_sorted,
    marker='o', linestyle='-', linewidth=1,
    markersize=3, markerfacecolor='white', markeredgewidth=0.8,
    color='black', label='# ms', zorder=5
)
ax0_2.plot(
    x_real, ds_counts_real_sorted,
    marker='s', linestyle='--', linewidth=1,
    markersize=3, markerfacecolor='white', markeredgewidth=0.8,
    color='black', label='# ds', zorder=5
)
ax0_2.plot(
    x_real, rpcs_real_sorted,
    marker='^', linestyle='--', linewidth=1,
    markersize=3, markerfacecolor='white', markeredgewidth=0.8,
    color='black', label='# rpcs', zorder=5
)

upper_lim_real = max(ms_counts_real_sorted.max(),
                     ds_counts_real_sorted.max(),
                     rpcs_real_sorted.max()) * 1.1
ax0_2.set_ylim(0, upper_lim_real)
ax0_2.tick_params(axis='y', labelsize=6)
ax0_2.set_ylabel('# ms / # ds / # rpcs', fontsize=6)
ax0_2.legend(
    loc='upper left',
    fontsize=6,
    ncol=3,
    handlelength=1.0,
    columnspacing=0.6,
    handletextpad=0.3,
)

axes[0].set_xticks(x_real)
axes[0].set_xticklabels(
    xtick_labels_real,
    rotation=20, ha="right", rotation_mode="anchor", fontsize=6
)

# =================== subplot 1: SYNTHETIC APPS ===================
bars1 = axes[1].bar(x_syn, total_syn_sorted, color=COLORS['total_syn'], width=bar_width)
axes[1].set_yscale("log")
axes[1].margins(y=0.2)  # 20% vertical padding
axes[1].bar_label(bars1, fmt="%.2fs", fontsize=5.5)
axes[1].set_title("Total (synthetic applications)", fontsize=6, pad=2)

ax1_2 = axes[1].twinx()

# overlay #ms / #ds / #rpcs for synthetic apps (log scale)
ax1_2.plot(
    x_syn, ms_counts_syn_sorted,
    marker='o', linestyle='-', linewidth=1,
    markersize=3, markerfacecolor='white', markeredgewidth=0.8,
    color='black', label='# ms', zorder=5
)
ax1_2.plot(
    x_syn, ds_counts_syn_sorted,
    marker='s', linestyle='--', linewidth=1,
    markersize=3, markerfacecolor='white', markeredgewidth=0.8,
    color='black', label='# ds', zorder=5
)
ax1_2.plot(
    x_syn, rpcs_syn_sorted,
    marker='^', linestyle='--', linewidth=1,
    markersize=3, markerfacecolor='white', markeredgewidth=0.8,
    color='black', label='# rpcs', zorder=5
)

upper_lim_syn = max(ms_counts_syn_sorted.max(),
                    ds_counts_syn_sorted.max(),
                    rpcs_syn_sorted.max()) * 2.2
ax1_2.tick_params(axis='y', labelsize=6)
ax1_2.set_ylabel('# ms / # ds / # rpcs', fontsize=6)
# no second legend here to avoid duplication

axes[1].set_xticks(x_syn)
axes[1].set_xticklabels(
    xtick_labels_syn,
    rotation=20, ha="right", rotation_mode="anchor", fontsize=6
)

# ------------ common formatting ------------
fig.text(0.02, 0.50, "Time (s)", va='center', rotation='vertical', fontsize=6)

for ax in axes:
    ax.tick_params(axis='x', length=2.5)
    ax.tick_params(axis='y', labelsize=5.5)

plt.tight_layout()
plt.subplots_adjust(left=0.12, hspace=0.40)

plt.savefig(OUTPUT_FILE1, bbox_inches='tight', pad_inches=0.05)
print(f"[INFO] saved plot to {OUTPUT_FILE1}")
plt.savefig(OUTPUT_FILE2, bbox_inches='tight', pad_inches=0.05)
print(f"[INFO] saved plot to {OUTPUT_FILE2}")
