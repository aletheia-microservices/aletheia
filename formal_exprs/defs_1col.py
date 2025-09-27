import matplotlib.pyplot as plt

# data
rows = [
    (r"$\mathrm{UNIQUE}_{T_{A}}(a_{i})$",
     "column $a_{i}$ is unique in table $T_{A}$."),
    (r"$\mathrm{PK}_{T_{A}}(a_{i})$",
     "column $a_{i}$ is a primary key in table $T_{A}$."),
    (r"$\mathrm{FK}\!\left(b_{j}^{T_{B}},\, a_{i}^{T_{A}}\right)$",
     "column $b_{j}$ in table $T_{B}$ is a foreign key referencing column $a_{i}$ in table $T_{A}$."),
    (r"$\mathrm{FK}\!\left(b_{j}^{T_{B}},\, a_{i}^{T_{A}},\, \mathrm{MANDATORY}\right)$",
     "a and b are created in the same request when the association is established."),
]

# fonts to look like your math
plt.rcParams["font.family"] = "cmr10"
plt.rcParams["mathtext.fontset"] = "cm"

fig, ax = plt.subplots(figsize=(10, 1.5))
ax.axis("off")

# build a table-like layout with two columns using fig.text for better wrapping
y = 0.80
line_gap = 0.19

for left, right in rows:
    fig.text(0, y, left, ha="left", va="top", fontsize=14)
    fig.text(0.30, y, right, ha="left", va="top", fontsize=13)
    y -= line_gap

fig.text(0, 0.99, "Notation", fontsize=14, weight="bold", va="top")
fig.text(0.30, 0.99, "Description",  fontsize=14, weight="bold", va="top")

plt.tight_layout()
plt.savefig("defs_1col.png", dpi=200, bbox_inches="tight")
