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

#plt.rcParams["text.usetex"] = True
plt.rcParams["font.family"] = "STIXGeneral"
plt.rcParams["mathtext.fontset"] = "stix"

fig = plt.figure(figsize=(1, 6))
ax = plt.gca()
ax.axis("off")

pad_title_expr = 0.015
pad_blocks     = 0.040
y = 0.98

def add_text_and_step(x, y, s, **kwargs):
    """add text at (x,y) in figure coords, measure its height, return new y"""
    t = fig.text(x, y, s, transform=fig.transFigure, **kwargs)
    fig.canvas.draw()
    renderer = fig.canvas.get_renderer()
    bbox_disp = t.get_window_extent(renderer=renderer)
    bbox_fig  = bbox_disp.transformed(fig.transFigure.inverted())
    return y - bbox_fig.height, t

for title, expr in exprs:
    # title
    y, _ = add_text_and_step(
        0, y, title,
        fontsize=16, va="top", ha="left", weight="medium"
    )
    y -= pad_title_expr

    # expression
    y, _ = add_text_and_step(
        0, y, expr,
        fontsize=16, va="top", ha="left", weight="medium"
    )
    y -= pad_blocks

plt.tight_layout()
plt.savefig("exprs_1col.png", dpi=200, bbox_inches="tight")
#plt.show()
