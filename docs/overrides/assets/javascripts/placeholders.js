// Marks the placeholders in code that readers replace with their own values, such as {app} or {spec}, so that
// page.css can make them stand out. A placeholder is a lowercase name in braces. Code blocks marked with the
// "warnings" class are skipped, since the warnings that Aletheia writes also use braces (e.g., database={cart_db}),
// and so are Go code and Mermaid diagrams, where braces are part of the code.
// A code block can also mark the values that fill in a placeholder, such as the name of the app in the tutorial,
// by listing them in its data-placeholders attribute (``` { .yaml data-placeholders="simpleshop" } ```).
// Only the display changes: the copy button still copies the plain text.
// document$ is Material's observable for each page load, including instant navigation.
// The script runs in its own scope, so its names do not clash with those of the other scripts.
(() => {
  const PLACEHOLDER = /\{[a-z][a-z0-9_]*\}/g;
  const SKIP = ".warnings, .language-go, .mermaid";

  // the placeholders and the values listed in the block's data-placeholders attribute (a space-separated list);
  // a value is marked only as a whole word, which may be followed by "_" (e.g., simpleshop in simpleshop_docker,
  // as {app} in {app}_{spec})
  const patternFor = (code) => {
    const values = code.closest("[data-placeholders]")?.dataset.placeholders.split(" ").filter(Boolean) ?? [];
    const escape = (value) => value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
    const sources = values.map((value) => `(?<![A-Za-z0-9])${escape(value)}(?![A-Za-z0-9])`);
    return new RegExp([PLACEHOLDER.source, ...sources].join("|"), "g");
  };

  // the syntax highlighter can split a placeholder into several tokens (e.g., "{", "app" and "}" in YAML), so
  // the placeholders are found in the text of the whole code element, and each text node gets its own spans
  const mark = (code) => {
    const ranges = [...code.textContent.matchAll(patternFor(code))].map((m) => [m.index, m.index + m[0].length]);
    if (!ranges.length) return;

    const walker = document.createTreeWalker(code, NodeFilter.SHOW_TEXT);
    const nodes = [];
    for (let node = walker.nextNode(); node; node = walker.nextNode()) nodes.push(node);

    let offset = 0;
    for (const node of nodes) {
      const text = node.data, start = offset, end = offset + text.length;
      offset = end;
      const parts = ranges
        .filter(([s, e]) => s < end && e > start)
        .map(([s, e]) => [Math.max(s, start) - start, Math.min(e, end) - start]);
      if (!parts.length) continue;

      const fragment = document.createDocumentFragment();
      let last = 0;
      for (const [s, e] of parts) {
        if (s > last) fragment.append(text.slice(last, s));
        const span = document.createElement("span");
        span.className = "placeholder";
        span.textContent = text.slice(s, e);
        fragment.append(span);
        last = e;
      }
      if (last < text.length) fragment.append(text.slice(last));
      node.replaceWith(fragment);
    }
  };

  document$.subscribe(() => {
    for (const code of document.querySelectorAll(".md-typeset code")) {
      if (code.closest(SKIP) || code.querySelector(".placeholder")) continue;
      mark(code);
    }
  });
})();
