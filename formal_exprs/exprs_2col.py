import matplotlib.pyplot as plt

expr1 = (
    r"$\exists a_{1}\in T_{A},\ b_{f}\in T_{B}:\ $"
    "\n"
    r"$\mathrm{FK}\!\left(b_{f}^{T_{B}},\, a_{1}^{T_{A}}\right)\ \wedge\ "
    r"d_{T_{A}}(a_{1}=x)\ \wedge\ \neg d_{T_{B}}(b_{f}=x)$"
)

expr2 = (
    r"$\forall x, y,\, 2 \leq f \leq N,\ \exists a_{1}\in T_{A},\ b_{1},\ b_{f}\in T_{B}:\ $"
    "\n"
    r"$\mathrm{FK}^{\neg\mathrm{TP}}\!\left(b_{f}^{T_{B}},\, a_{1}^{T_{A}}\right)\ \wedge\ "
    r"[d'_{T_{A}}(\ldots)]_{R2}\ \wedge\ "
    r"[\,r'_{T_{A}}(x)\ \wedge\ w'_{T_{B}}(y, x)\,]_{\mathrm{req_{1}}}$"
)

expr3 = (
    r"$\forall x, y,\, 2 \leq f \leq N,\ \exists a_{1}\in T_{A},\ b_{1},\ b_{f}\in T_{B}:\ $"
    "\n"
    r"$\mathrm{FK}^{\mathrm{TP}}\!\left(b_{f}^{T_{B}},\, a_{f}^{T_{A}}\right)\ \wedge\ "
    r"[\neg w_{T_{A}}(a_{1}=x)]_{\mathrm{req_{2}}}\ \wedge $"
    "\n"
    r"$([(b_{f}=x)=r_{T_{B}}(b_{1}=y)\ \wedge"
    r"\ r_{T_{A}}(a_{1}=x)]_{\mathrm{req_{2}}}\ \vee\ $"
    "\n"
    r"$[r_{T_{A}}(a_{1}=x)\ \wedge\ r_{T_{B}}(b_{f}=x)]_{\mathrm{req_{2}}})$"
)

expr4 = (
    r"$\forall k,\ \exists a_{1}\in T_{A},\ b_{1}\in T_{B}:\ $"
    "\n"
    r"$\mathrm{PK}_{T_{A}}(a_{1})\ \wedge\ \mathrm{PK}_{T_{B}}(b_{1})\ \wedge\ "
    r"\mathrm{FK}^{\mathrm{TP}}\!\left(b_{f}^{T_{B}},\, a_{f}^{T_{A}}\right)\ \wedge$"
    "\n"
    r"$[\neg\, w_{T_{A}}(a_{1}=k)]_{\mathrm{req_{2}}}\ \wedge\ "
    r"[\neg\, w_{T_{B}}(b_{1}=k)]_{\mathrm{req_{2}}}\ \wedge$"
    r"$[\,r'_{T_{A}}(k)\ \wedge\ r'_{T_{B}}(k)\,]_{R2}$"
)

expr5 = (
    r"$\forall x, y, v,\ 2 \leq v \leq N,\ \exists a_{1}, a_{v}\in T_{A},\ b_{1}\in T_{B}:\ $"
    "\n"
    r"$[(\mathrm{UNIQ}_{T_{A}}(a_{1}) \wedge w'_{T_{A}}(k))\ \vee\ "
    r"(\mathrm{UNIQ}_{T_{A}}(a_{v}) \wedge w'_{T_{A}}(k, v))]\ \wedge\ "
    r"w'_{T_{B}}(k)$"
)

exprs = [
    ("(P1) Referential Integrity - Absence of Cascade Delete", expr1),
    ("(P2) Referential Integrity - Concurrent Associations",   expr2),
    ("(P3) Referential Integrity - Uncoordinated Replication", expr3),
    ("(P4) Entity Integrity - Uncoordinated Replication",      expr4),
    ("(P5) Unicity - Conflict Resolution",                     expr5),
]

plt.rcParams["font.family"] = "STIXGeneral"
plt.rcParams["mathtext.fontset"] = "stix"

fig = plt.figure(figsize=(13, 3))   # taller to fit two columns nicely
ax = plt.gca()
ax.axis("off")

# layout params (figure coords)
col_x = [0.00, 0.50]   # left column x, right column x
pad_title_expr = 0.012 # gap between title and expr (inside a block)
pad_rows       = 0.080 # gap between rows (between blocks)
y = 0.97               # start near top

def add_text_get_height(x, y, s, **kwargs):
    """place text at (x,y) (figure coords), return its height in figure coords"""
    t = fig.text(x, y, s, transform=fig.transFigure, **kwargs)
    fig.canvas.draw()
    r = fig.canvas.get_renderer()
    bb = t.get_window_extent(renderer=r).transformed(fig.transFigure.inverted())
    return bb.height, t

def place_block(x, y, title, expr):
    """place one (title, expr) block at top y, return total height used"""
    # title
    h_title, _ = add_text_get_height(
        x, y, title, fontsize=16, va="top", ha="left"
    )
    # expr (placed below the title by its height + padding)
    y_expr = y - h_title - pad_title_expr
    h_expr, _ = add_text_get_height(
        x, y_expr, expr, fontsize=16, va="top", ha="left"
    )
    return (h_title + pad_title_expr + h_expr)

i = 0
while i < len(exprs):
    # left block
    titleL, exprL = exprs[i]
    h_left = place_block(col_x[0], y, titleL, exprL)

    # right block (if available)
    if i + 1 < len(exprs):
        titleR, exprR = exprs[i + 1]
        h_right = place_block(col_x[1], y, titleR, exprR)
    else:
        h_right = 0.0

    # move to next row by the taller of the two blocks, plus row padding
    y -= max(h_left, h_right) + pad_rows
    i += 2

plt.tight_layout()
plt.savefig("exprs_2col.png", dpi=200, bbox_inches="tight")
#plt.show()
