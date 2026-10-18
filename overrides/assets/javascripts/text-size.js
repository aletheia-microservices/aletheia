(() => {
  const KEY = "text-size";
  const MIN = 0.6, MAX = 1.2, STEP = 0.05, DEFAULT = 0.8;
  const root = document.documentElement;

  let size = DEFAULT;
  try { size = parseFloat(localStorage.getItem(KEY)) || DEFAULT; } catch {}

  const apply = (value) => {
    size = Math.round(Math.min(MAX, Math.max(MIN, value)) * 100) / 100;
    root.style.setProperty("--text-size", `${size}rem`);
    try { localStorage.setItem(KEY, size); } catch {}
  };

  apply(size);

  const header = document.querySelector(".md-header__inner");
  if (!header || header.querySelector(".text-size")) return;

  const group = document.createElement("div");
  group.className = "text-size";

  [
    ["A−", "Decrease text size", () => apply(size - STEP)],
    ["A",  "Reset text size",    () => apply(DEFAULT)],
    ["A+", "Increase text size", () => apply(size + STEP)],
  ].forEach(([label, title, onClick]) => {
    const btn = document.createElement("button");
    btn.textContent = label;
    btn.title = title;
    btn.setAttribute("aria-label", title);
    btn.addEventListener("click", onClick);
    group.appendChild(btn);
  });

  header.insertBefore(group, header.querySelector(".md-search"));
})();